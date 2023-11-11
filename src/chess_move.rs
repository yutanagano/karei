use crate::board::Coordinate;
use crate::piece::PieceType;

pub enum ChessMove {
    Standard{
        from: Coordinate,
        to: Coordinate
    },
    EnPassantCapture{
        from: Coordinate
    },
    PawnPromotion{
        from: Coordinate,
        promotion_to: PieceType
    },
    CastleKingside,
    CastleQueenside
}
