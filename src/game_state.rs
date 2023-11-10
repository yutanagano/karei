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

    pub fn try_from_fen(fen: &str) -> Result<Self, &str> {
        let mut fen_pieces = fen.split_whitespace();

        let piece_placement = match fen_pieces.next() {
            Some(s) => s,
            _ => return Err("Empty string.")
        };
        let mut board = [Square::new_empty(); 64];
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
        };

        let active_color = match fen_pieces.next() {
            Some(s) => match s {
                "w" => Color::White,
                "b" => Color::Black,
                _ => return Err("Malformed active color.")
            },
            _ => return Err("Incomplete FEN.")
        };

        let castling_info_string = match fen_pieces.next() {
            Some(s) => s,
            _ => return Err("Incomplete FEN.")
        };
        let mut castling_rights = CastlingRights::new_assuming_all_false();
        for char in castling_info_string.chars() {
            match char {
                'K' => castling_rights[Color::White].kingside = true,
                'Q' => castling_rights[Color::White].queenside = true,
                'k' => castling_rights[Color::Black].kingside = true,
                'q' => castling_rights[Color::Black].queenside = true,
                '-' => {},
                _ => return Err("Malformed castling rights.")
            }
        };

        let en_passant_square = match fen_pieces.next() {
            Some(s) => match s {
                "-" => None,
                coords => Some(Coordinate::try_from(coords).unwrap())
            },
            None => return Err("Incomplete FEN.")
        };
        if let Some(coord) = en_passant_square {
            board[usize::from(coord)].set_en_passant();
        };

        let plies_since_last_capture_or_pawn_advance: u8 = match fen_pieces.next() {
            Some(s) => s.parse().unwrap(),
            None => return Err("Incomplete FEN.")
        };

        let moves_played: u8 = match fen_pieces.next() {
            Some(s) => s.parse().unwrap(),
            None => return Err("Incomplete FEN.")
        };

        Ok(
            GameState {
                position: Position{ board, active_color, castling_rights },
                current_move_number: moves_played,
                num_plies_since_last_capture_or_pawn_advance: plies_since_last_capture_or_pawn_advance
            }
        )
    }

    pub fn print(&self) {
        self.position.print();
        println!("Current move number: {}", self.current_move_number);
        println!("Plies since last capture/pawn move: {}", self.num_plies_since_last_capture_or_pawn_advance);
    }
}


#[cfg(test)]
mod tests {
    use super::{GameState, Piece, PieceType, Color};

    #[test]
    fn instantiate_starting_game_state() {
        let gamestate = GameState::new();

        assert_eq!(
            gamestate.position.board[0].get_piece(),
            Some(Piece{ piece_type: PieceType::Rook, color: Color::Black })
        );
        assert_eq!(
            gamestate.position.board[49].get_piece(),
            Some(Piece{ piece_type: PieceType::Pawn, color: Color::White })
        );
        assert_eq!(
            gamestate.position.board[6].get_piece(),
            Some(Piece { piece_type: PieceType::Knight, color: Color::Black })
        );
        assert_eq!(
            gamestate.position.board[61].get_piece(),
            Some(Piece { piece_type: PieceType::Bishop, color: Color::White })
        );
        assert_eq!(
            gamestate.position.board[3].get_piece(),
            Some(Piece { piece_type: PieceType::Queen, color: Color::Black })
        );
        assert_eq!(
            gamestate.position.board[25].get_piece(),
            None
        );

        assert_eq!(gamestate.position.active_color, Color::White);

        assert!(gamestate.position.castling_rights[Color::White].kingside);
        assert!(gamestate.position.castling_rights[Color::White].queenside);
        assert!(gamestate.position.castling_rights[Color::Black].kingside);
        assert!(gamestate.position.castling_rights[Color::Black].queenside);

        assert_eq!(gamestate.current_move_number, 1);

        assert_eq!(gamestate.num_plies_since_last_capture_or_pawn_advance, 0);
    }


    #[test]
    fn read_arbitrary_fen() {
        let gamestate = match GameState::try_from_fen("4rk2/PR4pp/2R2p2/8/1P6/7P/5BK1/8 b - - 0 43") {
            Ok(gs) => gs,
            Err(msg) => panic!("{msg}")
        };

        assert_eq!(
            gamestate.position.board[4].get_piece(),
            Some(Piece{ piece_type: PieceType::Rook, color: Color::Black })
        );
        assert_eq!(
            gamestate.position.board[5].get_piece(),
            Some(Piece{ piece_type: PieceType::King, color: Color::Black })
        );
        assert_eq!(
            gamestate.position.board[31].get_piece(),
            None
        );
        assert_eq!(
            gamestate.position.board[33].get_piece(),
            Some(Piece{ piece_type: PieceType::Pawn, color: Color::White })
        );

        assert_eq!(gamestate.position.active_color, Color::Black);

        assert!(!gamestate.position.castling_rights[Color::White].kingside);
        assert!(!gamestate.position.castling_rights[Color::White].queenside);
        assert!(!gamestate.position.castling_rights[Color::Black].kingside);
        assert!(!gamestate.position.castling_rights[Color::Black].queenside);

        assert_eq!(gamestate.current_move_number, 43);

        assert_eq!(gamestate.num_plies_since_last_capture_or_pawn_advance, 0);
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
