use std::fmt;
use std::io::{self, Write};
use std::ops::{Index, IndexMut};


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

            board.squares[square_index] = Square {
                occupied_by: Some(Piece { piece_type, color }),
                en_passant_square: false
            };

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
        let mut castling_rights = CastlingRights{
            white_castling_rights: CastlingRightsForColor { kingside: false, queenside: false },
            black_castling_rights: CastlingRightsForColor { kingside: false, queenside: false }
        };
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
            board[coord].en_passant_square = true
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

        for square in self.board.squares.iter() {
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


#[derive(Copy, Clone)]
struct Board {
    squares: [Square; 64]
}

impl Board {
    fn new() -> Self {
        Board { squares: [Square::new_empty(); 64] }
    }
}

impl Index<Coordinate> for Board {
    type Output = Square;

    fn index(&self, index: Coordinate) -> &Self::Output {
        &self.squares[usize::from(index)]
    }
}

impl IndexMut<Coordinate> for Board {
    fn index_mut(&mut self, index: Coordinate) -> &mut Self::Output {
        &mut self.squares[usize::from(index)]
    }
}


struct Coordinate {
    index: usize
}

impl From<Coordinate> for usize {
    fn from(value: Coordinate) -> Self {
        value.index
    }
}

impl TryFrom<&str> for Coordinate {
    type Error = &'static str;

    fn try_from(value: &str) -> Result<Self, Self::Error> {
        let mut chars = value.chars();
        
        let file: usize = match chars.next() {
            Some(c) => match c.to_uppercase().next().unwrap() {
                'A' => 0,
                'B' => 1,
                'C' => 2,
                'D' => 3,
                'E' => 4,
                'F' => 5,
                'G' => 6,
                'H' => 7,
                _ => return Err("Unrecognized file.")
            },
            None => return Err("Can't convert empty string to coordinate.")
        };

        let rank = match chars.next() {
            Some(c) => match c {
                '1' => 7,
                '2' => 6,
                '3' => 5,
                '4' => 4,
                '5' => 3,
                '6' => 2,
                '7' => 1,
                '8' => 0,
                _ => return Err("Unrecognized rank.")
            }
            None => return Err("Missing rank.")
        };

        if !chars.next().is_none() {
            return Err("str too long.")
        }

        Ok(Coordinate { index: file + 8*rank })
    }
}



#[derive(Copy, Clone)]
struct Square {
    occupied_by: Option<Piece>,
    en_passant_square: bool
}

impl Square {
    fn new_empty() -> Self {
        Square { occupied_by: None, en_passant_square: false }
    }
}

impl fmt::Display for Square {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
         if self.occupied_by.is_none() {
            if self.en_passant_square {
                return write!(f, "*")
            } else {
                return write!(f, ".")
            }
        }

        let piece = self.occupied_by.unwrap();
        let letter = match piece.piece_type {
            PieceType::Pawn => 'P',
            PieceType::King => 'K',
            PieceType::Queen => 'Q',
            PieceType::Bishop => 'B',
            PieceType::Knight => 'K',
            PieceType::Rook => 'R'
        };
        if piece.color == Color::White {
            return write!(f, "{letter}")
        } else {
            return write!(f, "{}", letter.to_lowercase().next().unwrap())
        }
    }
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


#[derive(Copy, Clone, PartialEq)]
enum Color {
    White,
    Black
}

impl fmt::Display for Color {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let as_str = match self {
            Color::White => "white",
            Color::Black => "black"
        };
        write!(f, "{as_str}")
    }
}


struct CastlingRights {
    white_castling_rights: CastlingRightsForColor,
    black_castling_rights: CastlingRightsForColor
}

impl Index<Color> for CastlingRights {
    type Output = CastlingRightsForColor;

    fn index(&self, index: Color) -> &Self::Output {
        match index {
            Color::White => &self.white_castling_rights,
            Color::Black => &self.black_castling_rights
        }
    }
}

impl IndexMut<Color> for CastlingRights {
    fn index_mut(&mut self, index: Color) -> &mut Self::Output {
        match index {
            Color::White => &mut self.white_castling_rights,
            Color::Black => &mut self.black_castling_rights
        }
    }
}


struct CastlingRightsForColor {
    kingside: bool,
    queenside: bool
}

impl fmt::Display for CastlingRightsForColor {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        if self.kingside & self.queenside {
            write!(f, "kingside, queenside")
        } else if self.kingside {
            write!(f, "kingside")
        } else if self.queenside {
            write!(f, "queenside")
        } else {
            write!(f, "none")
        }
    }
}
