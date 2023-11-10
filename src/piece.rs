use crate::color::Color;


#[derive(Copy, Clone, Debug, PartialEq)]
pub struct Piece {
    pub piece_type: PieceType,
    pub color: Color
}

impl Piece {
    pub fn get_type(&self) -> PieceType {
        self.piece_type
    }

    pub fn get_color(&self) -> Color {
        self.color
    }
}


#[derive(Copy, Clone, Debug, PartialEq)]
pub enum PieceType {
    Pawn,
    King,
    Queen,
    Bishop,
    Knight,
    Rook
}
