use crate::board::{Board, Coordinate};
use crate::color::Color;
use crate::castling_rights::CastlingRights;
use crate::piece::{Piece, PieceType};
use std::io::{self, Write};


pub struct Position {
    board: Board,
    active_color: Color,
    moves_played: u8,
    plies_since_last_capture_or_pawn_advance: u8,
    castling_rights: CastlingRights
}

impl Position {
    pub fn new() -> Self {
        Position::from_fen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
    }

    pub fn from_fen(fen: &str) -> Self {
        let mut board = Board::new();
        let mut fen_pieces = fen.split_whitespace();

        let piece_placement = match fen_pieces.next() {
            Some(s) => s,
            _ => panic!("Invalid FEN: Can't read piece placements.")
        };

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
                    _ => panic!("Unrecognized piece character: {char}")
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
                _ => panic!("Unrecognized piece character: {char}")
            };

            let color = match char.is_ascii_uppercase() {
                true => Color::White,
                false => Color::Black
            };

            board[Coordinate::try_from(square_index).unwrap()].set_piece(
                Piece { piece_type, color }
            );

            square_index += 1;
        }

        let active_color = match fen_pieces.next() {
            Some(s) => match s {
                "w" => Color::White,
                "b" => Color::Black,
                _ => panic!("Invalid FEN: Unrecognized color designation: {}", s)
            },
            _ => panic!("Invalid FEN: Can't read active color.")
        };

        let castling_info_string = match fen_pieces.next() {
            Some(s) => s,
            _ => panic!("Invalid FEN: Can't read castling rights info.")
        };
        let mut castling_rights = CastlingRights::new_assuming_all_false();
        for char in castling_info_string.chars() {
            match char {
                'K' => castling_rights[Color::White].kingside = true,
                'Q' => castling_rights[Color::White].queenside = true,
                'k' => castling_rights[Color::Black].kingside = true,
                'q' => castling_rights[Color::Black].queenside = true,
                '-' => {},
                _ => panic!("Invalid FEN: Unrecognized castling right character: {}.", char)
            }
        };

        let en_passant_square = match fen_pieces.next() {
            Some(s) => match s {
                "-" => None,
                coords => Some(Coordinate::try_from(coords).unwrap())
            },
            None => panic!("Invalid FEN: Missing en passant info.")
        };
        if let Some(coord) = en_passant_square {
            board[coord].set_en_passant();
        };

        let plies_since_last_capture_or_pawn_advance: u8 = match fen_pieces.next() {
            Some(s) => s.parse().unwrap(),
            None => panic!("Invalid FEN: Missing ply info.")
        };

        let moves_played: u8 = match fen_pieces.next() {
            Some(s) => s.parse().unwrap(),
            None => panic!("Invalid FEN: Missing move count info.")
        };

        Position {
            board,
            active_color,
            moves_played,
            plies_since_last_capture_or_pawn_advance,
            castling_rights,
        }
    }

    pub fn print(&self) {
        let mut col_counter = 0;

        for square in self.board.iter_squares() {
            print!("{}", square);

            col_counter += 1;
            col_counter = col_counter % 8;

            if col_counter == 0 {
                print!("\n");
            }
        }
        io::stdout().flush().unwrap();

        println!("{} to move.", self.active_color);
        println!("White castling rights: {}.", self.castling_rights[Color::White]);
        println!("Black castling rights: {}.", self.castling_rights[Color::Black]);
    }
}



