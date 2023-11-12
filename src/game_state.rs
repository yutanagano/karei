use crate::board::{Coordinate, Square};
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

    pub fn try_from_fen(fen: &str) -> Result<Self, &'static str> {
        let mut fen_pieces = fen.split_whitespace();

        let mut board = [Square::new_empty(); 64];
        let active_color: Color;
        let castling_rights: CastlingRights;
        let current_move_number: u8;
        let num_plies_since_last_capture_or_pawn_advance: u8;

        let piece_placement = match fen_pieces.next() {
            Some(s) => s,
            None => return Err("Empty string.")
        };
        GameState::place_pieces(piece_placement, &mut board)?;

        let active_color_string = match fen_pieces.next() {
            Some(s) => s,
            None => return Err("Incomplete FEN.")
        };
        active_color = match active_color_string {
            "w" => Color::White,
            "b" => Color::Black,
            _ => return Err("Malformed active color.")
        };

        let castling_rights_string = match fen_pieces.next() {
            Some(s) => s,
            None => return Err("Incomplete FEN.")
        };
        castling_rights = GameState::parse_castling_rights(castling_rights_string)?;

        let en_passant_string = match fen_pieces.next() {
            Some(s) => s,
            None => return Err("Incomplete FEN.")
        };
        GameState::parse_en_passant_info(en_passant_string, &mut board)?;

        let ply_count_result = match fen_pieces.next() {
            Some(s) => s.parse::<u8>(),
            None => return Err("Incomplete FEN.")
        };
        num_plies_since_last_capture_or_pawn_advance = match ply_count_result {
            Ok(v) => v,
            Err(_) => return Err("Malformed ply count.")
        };

        let current_move_number_result = match fen_pieces.next() {
            Some(s) => s.parse::<u8>(),
            None => return Err("Incomplete FEN.")
        };
        current_move_number = match current_move_number_result {
            Ok(v) => v,
            Err(_) => return Err("Malformed move count.")
        };

        Ok(
            GameState {
                position: Position::new(board, active_color, castling_rights),
                current_move_number,
                num_plies_since_last_capture_or_pawn_advance
            }
        )
    }

    fn place_pieces(piece_placement: &str, board: &mut [Square; 64]) -> Result<(), &'static str> {
        let mut square_index: usize = 0;

        for char in piece_placement.chars() {
            if !char.is_ascii_alphabetic() {
                match char {
                    '1' => square_index += 1,
                    '2' => square_index += 2,
                    '3' => square_index += 3,
                    '4' => square_index += 4,
                    '5' => square_index += 5,
                    '6' => square_index += 6,
                    '7' => square_index += 7,
                    '8' => square_index += 8,
                    '/' => {},
                    _ => return Err("Malformed piece placements.")
                };
                continue;
            };

            let piece_type = match char.to_ascii_uppercase() {
                'P' => PieceType::Pawn,
                'K' => PieceType::King,
                'Q' => PieceType::Queen,
                'B' => PieceType::Bishop,
                'N' => PieceType::Knight,
                'R' => PieceType::Rook,
                _ => return Err("Malformed piece placements.")
            };

            let color = match char.is_ascii_uppercase() {
                true => Color::White,
                false => Color::Black
            };

            board[square_index].set_piece(Piece { piece_type, color });

            square_index += 1;
        }

        if square_index != 64 {
            return Err("Incomplete piece placements.")
        }

        Ok(())
    }

    fn parse_castling_rights(castling_rights_string: &str) -> Result<CastlingRights, &'static str> {
        let mut castling_rights = CastlingRights::new_assuming_all_false();

        for char in castling_rights_string.chars() {
            match char {
                'K' => castling_rights[Color::White].kingside = true,
                'Q' => castling_rights[Color::White].queenside = true,
                'k' => castling_rights[Color::Black].kingside = true,
                'q' => castling_rights[Color::Black].queenside = true,
                '-' => {},
                _ => return Err("Malformed castling rights.")
            }
        }

        Ok(castling_rights)
    }

    fn parse_en_passant_info(en_passant_string: &str, board: &mut [Square; 64]) -> Result<(), &'static str> {
        let maybe_coord = match en_passant_string {
            "-" => None,
            coord_str => Some(Coordinate::try_from(coord_str))
        };

        if let Some(coord_result) = maybe_coord {
            match coord_result {
                Ok(coord) => board[coord as usize].set_en_passant(),
                Err(_) => return Err("Malformed en passant square.")
            }
        };

        Ok(())
    }

    pub fn print(&self) {
        self.position.print();
        println!("Current move number: {}", self.current_move_number);
        println!("Plies since last capture/pawn move: {}", self.num_plies_since_last_capture_or_pawn_advance);
    }
}


#[cfg(test)]
mod tests {
    use super::{GameState, Coordinate, Square, Piece, PieceType, Color};

    #[test]
    fn test_place_pieces() {
        let mut board = [Square::new_empty(); 64];

        let _ = GameState::place_pieces("4rk2/PR4pp/2R2p2/8/1P6/7P/5BK1/8", &mut board).unwrap();

        assert_eq!(
            board[Coordinate::E8 as usize].get_piece(),
            Some(Piece{ piece_type: PieceType::Rook, color: Color::Black })
        );
        assert_eq!(
            board[Coordinate::F8 as usize].get_piece(),
            Some(Piece{ piece_type: PieceType::King, color: Color::Black })
        );
        assert_eq!(
            board[Coordinate::H5 as usize].get_piece(),
            None
        );
        assert_eq!(
            board[Coordinate::B4 as usize].get_piece(),
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
            Err(msg) => assert_eq!(msg, "Empty string."),
            Ok(_) => panic!("GameState from empty FEN did not return Err.")
        }
    }

    #[test]
    fn try_from_fen_with_bad_piece_placements() {
        match GameState::try_from_fen("foobarbaz w - - 5 7") {
            Err(msg) => assert_eq!(msg, "Malformed piece placements."),
            Ok(_) => panic!("GameState from FEN with bad piece placements did not return Err.")
        }
    }
}
