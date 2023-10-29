use crate::color::Color;
use std::fmt;
use std::ops::{Index, IndexMut};


pub struct CastlingRights {
    white_castling_rights: CastlingRightsForColor,
    black_castling_rights: CastlingRightsForColor
}

impl CastlingRights {
    pub fn new_assuming_all_false() -> Self {
        CastlingRights {
            white_castling_rights: CastlingRightsForColor { kingside: false, queenside: false },
            black_castling_rights: CastlingRightsForColor { kingside: false, queenside: false }
        }
    }
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


pub struct CastlingRightsForColor {
    pub kingside: bool,
    pub queenside: bool
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
