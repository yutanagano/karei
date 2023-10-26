use std::ops::{Index, IndexMut};

struct Position {
    board: Board,
    active_color: Color,
    moves_played: u8,
    plies_since_last_capture_or_pawn_advance: u8,
    white_castling_rights: CastlingRights,
    black_castling_rights: CastlingRights
}

impl Position {
    fn new() -> Self {
        Position::from_fen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
    }

    fn from_fen(fen: &str) -> Self {
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

            board[square_index] = Square {
                occupied_by: Some(Piece { piece_type, color }),
                en_passant_square: false
            };
        }

        Position {
            board,
            active_color: Color::White,
            moves_played: 1,
            plies_since_last_capture_or_pawn_advance: 0,
            white_castling_rights: CastlingRights { kingside: true, queenside: true },
            black_castling_rights: CastlingRights { kingside: true, queenside: true }
        }
    }
}

struct Board {
    squares: [Square; 64]
}

impl Board {
    fn new() -> Self {
        Board {
            squares: [Square { occupied_by: None, en_passant_square: false }; 64]
        }
    }
}

impl Index<usize> for Board {
    type Output = Square;

    fn index(&self, index: usize) -> &Self::Output {
        &self.squares[index]
    }
}

impl IndexMut<usize> for Board {
    fn index_mut(&mut self, index: usize) -> &mut Self::Output {
        &mut self.squares[index]
    }
}

#[derive(Copy, Clone)]
struct Square {
    occupied_by: Option<Piece>,
    en_passant_square: bool
}

#[derive(Copy, Clone)]
struct Piece {
    piece_type: PieceType,
    color: Color
}

#[derive(Copy, Clone)]
enum PieceType {
    Pawn,
    King,
    Queen,
    Bishop,
    Knight,
    Rook
}

#[derive(Copy, Clone)]
enum Color {
    White,
    Black
}

struct CastlingRights {
    kingside: bool,
    queenside: bool
}

fn main() {
    println!("Hello, world!");
    let _position = Position::new();
}
