use crate::color::Color;
use crate::piece::{Piece, PieceType};
use std::fmt;
use std::ops::{Index, IndexMut};


#[derive(Copy, Clone)]
pub struct Board {
    squares: [Square; 64]
}

impl Board {
    pub fn new() -> Self {
        Board { squares: [Square::new_empty(); 64] }
    }

    pub fn iter_squares(&self) -> std::slice::Iter<Square> {
        self.squares.iter()
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


#[derive(Copy, Clone)]
pub struct Square {
    occupied_by: Option<Piece>,
    en_passant_square: bool
}

impl Square {
    fn new_empty() -> Self {
        Square { occupied_by: None, en_passant_square: false }
    }

    pub fn set_piece(&mut self, piece: Piece) {
        self.occupied_by = Some(piece);
    }

    pub fn clear_piece(&mut self) {
        self.occupied_by = None;
    }

    pub fn set_en_passant(&mut self) {
        self.en_passant_square = true;
    }
}

impl fmt::Display for Square {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
         if self.occupied_by.is_none() {
            if self.en_passant_square {
                return write!(f, "*")
            } else {
                return write!(f, " ")
            }
        }

        let piece = self.occupied_by.unwrap();
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
