use crate::board::{Board, Coordinate, File, Rank, BoardError};
use crate::castling_rights::CastlingRights;
use crate::color::Color;
use crate::piece::{Piece, PieceType};
use crate::position::Position;


pub struct GameState {
    position: Position,
    current_move_number: u8,
    num_plies_since_last_capture_or_pawn_advance: u8
}

impl GameState {
    pub fn new() -> Self {
        GameState::try_from_fen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1").unwrap()
    }

    pub fn try_from_fen(fen: &str) -> Result<Self, FenError> {
        let mut fen_pieces = fen.split_whitespace();

        let mut board = Board::new_empty();
        let active_color: Color;
        let castling_rights: CastlingRights;
        let current_move_number: u8;
        let num_plies_since_last_capture_or_pawn_advance: u8;

        let piece_placement = fen_pieces.next().ok_or(FenError::EmptyString)?;
        GameState::place_pieces(piece_placement, &mut board)?;

        let active_color_string = fen_pieces.next().ok_or(FenError::Incomplete)?;
        active_color = match active_color_string {
            "w" => Color::White,
            "b" => Color::Black,
            _ => return Err(FenError::Malformed)
        };

        let castling_rights_string = fen_pieces.next().ok_or(FenError::Incomplete)?;
        castling_rights = GameState::parse_castling_rights(castling_rights_string)?;

        let en_passant_string = fen_pieces.next().ok_or(FenError::Incomplete)?;
        GameState::parse_en_passant_info(en_passant_string, &mut board)?;

        num_plies_since_last_capture_or_pawn_advance = match fen_pieces.next() {
            Some(s) => s.parse::<u8>().map_err(|_| FenError::Malformed)?,
            None => return Err(FenError::Incomplete)
        };

        current_move_number = match fen_pieces.next() {
            Some(s) => s.parse::<u8>().map_err(|_| FenError::Malformed)?,
            None => return Err(FenError::Incomplete)
        };

        Ok(
            GameState {
                position: Position::new(board, active_color, castling_rights),
                current_move_number,
                num_plies_since_last_capture_or_pawn_advance
            }
        )
    }

    fn place_pieces(piece_placement: &str, board: &mut Board) -> Result<(), FenError> {
        let mut square_index = 0;

        for char in piece_placement.chars() {
            let move_forward_by;

            if char.is_ascii_digit() {
                move_forward_by = char.to_digit(10).unwrap() as i8;
            } else if char.is_ascii_alphabetic() {
                move_forward_by = 1;

                let piece_type = match char.to_ascii_uppercase() {
                    'P' => PieceType::Pawn,
                    'K' => PieceType::King,
                    'Q' => PieceType::Queen,
                    'B' => PieceType::Bishop,
                    'N' => PieceType::Knight,
                    'R' => PieceType::Rook,
                    _ => return Err(FenError::Malformed)
                };

                let color = match char.is_ascii_uppercase() {
                    true => Color::White,
                    false => Color::Black
                };

                let current_file_index = square_index % 8;
                let current_rank_index = 7 - (square_index / 8);
                let current_coordinate = Coordinate::new(
                    File::try_from(current_file_index)?,
                    Rank::try_from(current_rank_index)?
                );

                board[current_coordinate].set_piece(Piece { piece_type, color });
            } else if char == '/' {
                move_forward_by = 0;
            } else {
                return Err(FenError::Malformed)
            }
            
            square_index += move_forward_by;
        }

        if square_index != 64 {
            return Err(FenError::Malformed)
        }

        Ok(())
    }

    fn parse_castling_rights(castling_rights_string: &str) -> Result<CastlingRights, FenError> {
        let mut castling_rights = CastlingRights::new_assuming_all_false();

        for char in castling_rights_string.chars() {
            match char {
                'K' => castling_rights[Color::White].kingside = true,
                'Q' => castling_rights[Color::White].queenside = true,
                'k' => castling_rights[Color::Black].kingside = true,
                'q' => castling_rights[Color::Black].queenside = true,
                '-' => {},
                _ => return Err(FenError::Malformed)
            }
        }

        Ok(castling_rights)
    }

    fn parse_en_passant_info(en_passant_string: &str, board: &mut Board) -> Result<(), FenError> {
        let maybe_coord = match en_passant_string {
            "-" => None,
            coord_str => {
                let mut chars = coord_str.chars();
        
                let file = match chars.next() {
                    Some(c) => match c.to_uppercase().next().unwrap() {
                        'A' => File::A,
                        'B' => File::B,
                        'C' => File::C,
                        'D' => File::D,
                        'E' => File::E,
                        'F' => File::F,
                        'G' => File::G,
                        'H' => File::H,
                        _ => return Err(FenError::Malformed)
                    },
                    None => return Err(FenError::Malformed)
                };

                let rank = match chars.next() {
                    Some(c) => match c {
                        '1' => Rank::First,
                        '2' => Rank::Second,
                        '3' => Rank::Third,
                        '4' => Rank::Fourth,
                        '5' => Rank::Fifth,
                        '6' => Rank::Sixth,
                        '7' => Rank::Seventh,
                        '8' => Rank::Eighth,
                        _ => return Err(FenError::Malformed)
                    }
                    None => return Err(FenError::Malformed)
                };

                if !chars.next().is_none() {
                    return Err(FenError::Malformed)
                }

                Some(Coordinate{ file, rank })
            }
        };

        if let Some(coord) = maybe_coord {
            board[coord].set_en_passant();
        };

        Ok(())
    }

    pub fn print(&self) {
        self.position.print();
        println!("Current move number: {}", self.current_move_number);
        println!("Plies since last capture/pawn move: {}", self.num_plies_since_last_capture_or_pawn_advance);
    }
}


#[derive(Debug, PartialEq)]
pub enum FenError {
    EmptyString,
    Incomplete,
    Malformed,
}

impl From<BoardError> for FenError {
    fn from(_: BoardError) -> Self {
        FenError::Malformed
    }
}


#[cfg(test)]
mod tests {
    use super::{GameState, Board, Coordinate, File, Rank, Piece, PieceType, Color, FenError};

    #[test]
    fn test_place_pieces() {
        let mut board = Board::new_empty();

        let _ = GameState::place_pieces("4rk2/PR4pp/2R2p2/8/1P6/7P/5BK1/8", &mut board).unwrap();

        assert_eq!(
            board[Coordinate::new(File::E, Rank::Eighth)].get_piece(),
            Some(Piece{ piece_type: PieceType::Rook, color: Color::Black })
        );
        assert_eq!(
            board[Coordinate::new(File::F, Rank::Eighth)].get_piece(),
            Some(Piece{ piece_type: PieceType::King, color: Color::Black })
        );
        assert_eq!(
            board[Coordinate::new(File::H, Rank::Fifth)].get_piece(),
            None
        );
        assert_eq!(
            board[Coordinate::new(File::B, Rank::Fourth)].get_piece(),
            Some(Piece{ piece_type: PieceType::Pawn, color: Color::White })
        );
    }

    #[test]
    fn test_parse_castling_rights() {
        let castling_rights = GameState::parse_castling_rights("KQkq").unwrap();

        assert!(castling_rights[Color::White].kingside);
        assert!(castling_rights[Color::White].queenside);
        assert!(castling_rights[Color::Black].kingside);
        assert!(castling_rights[Color::Black].queenside);

        let castling_rights = GameState::parse_castling_rights("Qk").unwrap();

        assert!(!castling_rights[Color::White].kingside);
        assert!(castling_rights[Color::White].queenside);
        assert!(castling_rights[Color::Black].kingside);
        assert!(!castling_rights[Color::Black].queenside);

        let castling_rights = GameState::parse_castling_rights("-").unwrap();

        assert!(!castling_rights[Color::White].kingside);
        assert!(!castling_rights[Color::White].queenside);
        assert!(!castling_rights[Color::Black].kingside);
        assert!(!castling_rights[Color::Black].queenside);
    }

    #[test]
    fn try_from_empty_fen() {
        match GameState::try_from_fen("") {
            Err(msg) => assert_eq!(msg, FenError::EmptyString),
            Ok(_) => panic!("GameState from empty FEN did not return Err.")
        }
    }

    #[test]
    fn try_from_fen_with_bad_piece_placements() {
        match GameState::try_from_fen("foobarbaz w - - 5 7") {
            Err(msg) => assert_eq!(msg, FenError::Malformed),
            Ok(_) => panic!("GameState from FEN with bad piece placements did not return Err.")
        }
    }
}
