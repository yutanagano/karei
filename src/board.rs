use crate::color::Color;
use crate::piece::{Piece, PieceType};
use std::fmt;


pub struct Coordinate {
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

impl TryFrom<usize> for Coordinate {
    type Error = &'static str;

    fn try_from(value: usize) -> Result<Self, Self::Error> {
        if value > 63 {
            return Err("Index out of bounds.")
        }

        Ok(Coordinate { index: value })
    }
}


#[derive(Copy, Clone, Debug)]
pub struct Square {
    piece: Option<Piece>,
    is_en_passant_square: bool,
    is_attacked_by_white: bool,
    is_attacked_by_black: bool
}

impl Square {
    pub fn new_empty() -> Self {
        Square { piece: None, is_en_passant_square: false, is_attacked_by_white: false, is_attacked_by_black: false }
    }

    pub fn get_piece(&self) -> Option<Piece> {
        self.piece
    }

    pub fn set_piece(&mut self, piece: Piece) {
        self.piece = Some(piece);
    }

    pub fn clear_piece(&mut self) {
        self.piece = None;
    }

    pub fn set_en_passant(&mut self) {
        self.is_en_passant_square = true;
    }

    pub fn is_empty(&self) -> bool {
        self.piece == None
    }

    pub fn is_attacked_by(&self, color: Color) -> bool {
        match color {
            Color::White => self.is_attacked_by_white,
            Color::Black => self.is_attacked_by_black
        }
    }
}

impl fmt::Display for Square {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
         if self.piece.is_none() {
            if self.is_en_passant_square {
                return write!(f, "*")
            } else {
                return write!(f, " ")
            }
        }

        let piece = self.piece.unwrap();
        let letter = match piece.get_type() {
            PieceType::Pawn => 'P',
            PieceType::King => 'K',
            PieceType::Queen => 'Q',
            PieceType::Bishop => 'B',
            PieceType::Knight => 'N',
            PieceType::Rook => 'R'
        };
        if piece.get_color() == Color::White {
            return write!(f, "{letter}")
        } else {
            return write!(f, "{}", letter.to_lowercase().next().unwrap())
        }
    }
}
