use crate::color::Color;
use crate::piece::{Piece, PieceType};
use std::fmt;
use std::ops::Add;


#[derive(Clone, Copy)]
pub enum Coordinate {
    A8, B8, C8, D8, E8, F8, G8, H8,
    A7, B7, C7, D7, E7, F7, G7, H7,
    A6, B6, C6, D6, E6, F6, G6, H6,
    A5, B5, C5, D5, E5, F5, G5, H5,
    A4, B4, C4, D4, E4, F4, G4, H4,
    A3, B3, C3, D3, E3, F3, G3, H3,
    A2, B2, C2, D2, E2, F2, G2, H2,
    A1, B1, C1, D1, E1, F1, G1, H1
}

impl Coordinate {
    pub fn get_rank(&self) -> Rank {
        match *self as u8 {
            (0..=7) => Rank::Eighth,
            (8..=15) => Rank::Seventh,
            (16..=23) => Rank::Sixth,
            (24..=31) => Rank::Fifth,
            (32..=39) => Rank::Fourth,
            (40..=47) => Rank::Third,
            (48..=55) => Rank::Second,
            (56..=63) => Rank::First,
            _ => panic!("Something went terribly wrong.")
        }
    }

    pub fn get_file(&self) -> File {
        match *self as u8 % 8 {
            0 => File::A,
            1 => File::B,
            2 => File::C,
            3 => File::D,
            4 => File::E,
            5 => File::F,
            6 => File::G,
            7 => File::H,
            _ => panic!("Something went terribly wrong.")
        }
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

        Coordinate::try_from(file + 8*rank)
    }
}

impl TryFrom<usize> for Coordinate {
    type Error = &'static str;

    fn try_from(value: usize) -> Result<Self, Self::Error> {
        match value {
            0 => Ok(Coordinate::A8),
            1 => Ok(Coordinate::B8),
            2 => Ok(Coordinate::C8),
            3 => Ok(Coordinate::D8),
            4 => Ok(Coordinate::E8),
            5 => Ok(Coordinate::F8),
            6 => Ok(Coordinate::G8),
            7 => Ok(Coordinate::H8),
            8 => Ok(Coordinate::A7),
            9 => Ok(Coordinate::B7),
            10 => Ok(Coordinate::C7),
            11 => Ok(Coordinate::D7),
            12 => Ok(Coordinate::E7),
            13 => Ok(Coordinate::F7),
            14 => Ok(Coordinate::G7),
            15 => Ok(Coordinate::H7),
            16 => Ok(Coordinate::A6),
            17 => Ok(Coordinate::B6),
            18 => Ok(Coordinate::C6),
            19 => Ok(Coordinate::D6),
            20 => Ok(Coordinate::E6),
            21 => Ok(Coordinate::F6),
            22 => Ok(Coordinate::G6),
            23 => Ok(Coordinate::H6),
            24 => Ok(Coordinate::A5),
            25 => Ok(Coordinate::B5),
            26 => Ok(Coordinate::C5),
            27 => Ok(Coordinate::D5),
            28 => Ok(Coordinate::E5),
            29 => Ok(Coordinate::F5),
            30 => Ok(Coordinate::G5),
            31 => Ok(Coordinate::H5),
            32 => Ok(Coordinate::A4),
            33 => Ok(Coordinate::B4),
            34 => Ok(Coordinate::C4),
            35 => Ok(Coordinate::D4),
            36 => Ok(Coordinate::E4),
            37 => Ok(Coordinate::F4),
            38 => Ok(Coordinate::G4),
            39 => Ok(Coordinate::H4),
            40 => Ok(Coordinate::A3),
            41 => Ok(Coordinate::B3),
            42 => Ok(Coordinate::C3),
            43 => Ok(Coordinate::D3),
            44 => Ok(Coordinate::E3),
            45 => Ok(Coordinate::F3),
            46 => Ok(Coordinate::G3),
            47 => Ok(Coordinate::H3),
            48 => Ok(Coordinate::A2),
            49 => Ok(Coordinate::B2),
            50 => Ok(Coordinate::C2),
            51 => Ok(Coordinate::D2),
            52 => Ok(Coordinate::E2),
            53 => Ok(Coordinate::F2),
            54 => Ok(Coordinate::G2),
            55 => Ok(Coordinate::H2),
            56 => Ok(Coordinate::A1),
            57 => Ok(Coordinate::B1),
            58 => Ok(Coordinate::C1),
            59 => Ok(Coordinate::D1),
            60 => Ok(Coordinate::E1),
            61 => Ok(Coordinate::F1),
            62 => Ok(Coordinate::G1),
            63 => Ok(Coordinate::H1),
            _ => Err("Index out of bounds.")
        }
    }
}

impl Add<Direction> for Coordinate {
    type Output = Result<Self, &'static str>;

    fn add(self, rhs: Direction) -> Self::Output {
        let new_index = (self as i8) + (i8::from(rhs));
        Coordinate::try_from(new_index as usize)
    }
}


#[derive(PartialEq)]
pub enum Rank {
    First,
    Second,
    Third,
    Fourth,
    Fifth,
    Sixth,
    Seventh,
    Eighth
}


#[derive(PartialEq)]
pub enum File {
    A, B, C, D, E, F, G, H
}


#[derive(Clone, Copy)]
pub enum Direction {
    Up,
    Down,
    Left,
    Right,
    UpperLeft,
    UpperRight,
    LowerLeft,
    LowerRight,
    Compound(i8)
}

impl Add for Direction {
    type Output = Self;

    fn add(self, rhs: Self) -> Self::Output {
        Direction::Compound(i8::from(self) + i8::from(rhs))
    }
}

impl From<Direction> for i8 {
    fn from(value: Direction) -> Self {
        match value {
            Direction::Up => -8,
            Direction::Down => 8,
            Direction::Left => -1,
            Direction::Right => 1,
            Direction::UpperLeft => -9,
            Direction::UpperRight => -7,
            Direction::LowerLeft => 7,
            Direction::LowerRight => 9,
            Direction::Compound(val) => val
        }
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

    pub fn has_piece_of_color(&self, color: Color) -> bool {
        if let Some(piece) = self.piece {
            piece.color == color
        } else {
            false
        }
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
