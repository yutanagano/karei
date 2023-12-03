use crate::color::Color;
use crate::piece::{Piece, PieceType};
use std::fmt;
use std::ops::{Index, IndexMut};


pub struct Board {
    squares: [[Square; 8]; 8]
}

impl Board {
    pub fn new_empty() -> Self {
        Board {
            squares: [[Square::new_empty(); 8]; 8]
        }
    }
}

impl Index<Coordinate> for Board {
    type Output = Square;

    fn index(&self, index: Coordinate) -> &Self::Output {
        &self.squares[index.file as usize][index.rank as usize]
    }
}

impl IndexMut<Coordinate> for Board {
    fn index_mut(&mut self, index: Coordinate) -> &mut Self::Output {
        &mut self.squares[index.file as usize][index.rank as usize]
    }
}


#[derive(Clone, Copy, Debug, PartialEq)]
pub struct Coordinate {
    file: File,
    rank: Rank
}

impl Coordinate {
    pub fn new(file: File, rank: Rank) -> Self {
        Coordinate { file, rank }
    }

    pub fn get_file(&self) -> File {
        return self.file;
    }

    pub fn get_rank(&self) -> Rank {
        return self.rank;
    }

    pub fn try_moving(self, direction: Direction) -> Result<Self, BoardError> {
        let file = self.file.try_moving(direction.delta_file)?;
        let rank = self.rank.try_moving(direction.delta_rank)?;
        Ok(Coordinate { file, rank })
    }
}


#[derive(Clone, Copy, Debug, PartialEq)]
pub enum File {
    A, B, C, D, E, F, G, H
}

impl File {
    pub fn try_moving(self, delta: i8) -> Result<Self, BoardError> {
        File::try_from(self as i8 + delta)
    }
}

impl TryFrom<i8> for File {
    type Error = BoardError;

    fn try_from(value: i8) -> Result<Self, Self::Error> {
        match value {
            0 => Ok(File::A),
            1 => Ok(File::B),
            2 => Ok(File::C),
            3 => Ok(File::D),
            4 => Ok(File::E),
            5 => Ok(File::F),
            6 => Ok(File::G),
            7 => Ok(File::H),
            _ => Err(BoardError::FileOutOfBounds)
        }
    }
}


#[derive(Clone, Copy, Debug, PartialEq)]
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

impl Rank {
    pub fn try_moving(self, delta: i8) -> Result<Self, BoardError> {
        Rank::try_from(self as i8 + delta)
    }
}

impl TryFrom<i8> for Rank {
    type Error = BoardError;

    fn try_from(value: i8) -> Result<Self, Self::Error> {
        match value {
            0 => Ok(Rank::First),
            1 => Ok(Rank::Second),
            2 => Ok(Rank::Third),
            3 => Ok(Rank::Fourth),
            4 => Ok(Rank::Fifth),
            5 => Ok(Rank::Sixth),
            6 => Ok(Rank::Seventh),
            7 => Ok(Rank::Eighth),
            _ => Err(BoardError::RankOutOfBounds)
        }
    }
}


#[derive(Clone, Copy)]
pub struct Direction {
    delta_file: i8,
    delta_rank: i8
}

impl Direction {
    pub fn new(delta_file: i8, delta_rank: i8) -> Self {
        Direction { delta_file, delta_rank }
    }

    pub fn up() -> Self {
        Direction { delta_file: 0, delta_rank: 1 }
    }

    pub fn down() -> Self {
        Direction { delta_file: 0, delta_rank: -1 }
    }

    pub fn left() -> Self {
        Direction { delta_file: -1, delta_rank: 0 }
    }

    pub fn right() -> Self {
        Direction { delta_file: 1, delta_rank: 0 }
    }

    pub fn upper_left() -> Self {
        Direction { delta_file: -1, delta_rank: 1 }
    }

    pub fn upper_right() -> Self {
        Direction { delta_file: 1, delta_rank: 1 }
    }

    pub fn lower_left() -> Self {
        Direction { delta_file: -1, delta_rank: -1 }
    }

    pub fn lower_right() -> Self {
        Direction { delta_file: 1, delta_rank: -1 }
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


#[derive(Debug)]
pub enum BoardError {
    FileOutOfBounds,
    RankOutOfBounds
}


#[cfg(test)]
mod tests {
    use super::{Coordinate, File, Rank, Direction};

    #[test]
    fn test_moving_coordinate() {
        let mut coordinate = Coordinate::new(File::A, Rank::First);
        coordinate = coordinate.try_moving(Direction::new(1, 0)).unwrap();
        assert_eq!(coordinate, Coordinate::new(File::B, Rank::First));
    }
}
