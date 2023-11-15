use crate::board::{Coordinate, Direction, Square, Rank, Board};
use crate::piece::PieceType;
use crate::chess_move::ChessMove;
use crate::color::Color;
use crate::castling_rights::CastlingRights;
use std::io::{self, Write};
use std::ops::{Index, IndexMut};


pub struct Position {
    board: Board,
    active_color: Color,
    castling_rights: CastlingRights
}

impl Position {
    pub fn new(board: Board, active_color: Color, castling_rights: CastlingRights) -> Self {
        Position { board, active_color, castling_rights }
    }

    pub fn get_possible_moves(&self) -> Vec<ChessMove> {
        let mut moves = Vec::new();

//        for (index, square) in self.board.iter().enumerate() {
//            if square.is_empty() {
//                continue;
//            };
//
//            let piece = square.get_piece().unwrap();
//
//            if piece.color == self.active_color.get_opposite() {
//                continue;
//            };
//
//            let from_coordinate = Coordinate::try_from(index).unwrap();
//
//            let pseudo_legal_moves = match piece.piece_type {
//                PieceType::Pawn => self.get_possible_pawn_moves(from_coordinate),
//                PieceType::King => self.get_possible_king_moves(from_coordinate),
//                PieceType::Queen => self.get_possible_queen_moves(from_coordinate),
//                PieceType::Bishop => self.get_possible_bishop_moves(from_coordinate),
//                PieceType::Knight => self.get_possible_knight_moves(from_coordinate),
//                PieceType::Rook => self.get_possible_rook_moves(from_coordinate)
//            };
//
//            moves.extend(pseudo_legal_moves);
//        }

        moves
    }

//    fn get_possible_pawn_moves(&self, from_coordinate: Coordinate) -> Vec<ChessMove> {
//        let mut moves = Vec::new();
//        let coordinate_above = (from_coordinate + Direction::Up).unwrap();
//        
//        if from_coordinate.get_rank() == Rank::Seventh {
//            for direction in [Direction::UpperLeft, Direction::UpperRight].into_iter() {
//                let to_coordinate = (from_coordinate + direction).unwrap();
//
//                if self.board[to_coordinate as usize].has_piece_of_color(self.active_color.get_opposite()) {
//                    moves.push(ChessMove::PawnPromotion { from: from_coordinate, to: to_coordinate, promotion_to: PieceType::Queen });
//                    moves.push(ChessMove::PawnPromotion { from: from_coordinate, to: to_coordinate, promotion_to: PieceType::Rook });
//                    moves.push(ChessMove::PawnPromotion { from: from_coordinate, to: to_coordinate, promotion_to: PieceType::Bishop });
//                    moves.push(ChessMove::PawnPromotion { from: from_coordinate, to: to_coordinate, promotion_to: PieceType::Knight });
//                }
//            }
//
//            if self.board[coordinate_above as usize].is_empty() {
//                moves.push(ChessMove::PawnPromotion { from: from_coordinate, to: coordinate_above, promotion_to: PieceType::Queen });
//                moves.push(ChessMove::PawnPromotion { from: from_coordinate, to: coordinate_above, promotion_to: PieceType::Rook });
//                moves.push(ChessMove::PawnPromotion { from: from_coordinate, to: coordinate_above, promotion_to: PieceType::Bishop });
//                moves.push(ChessMove::PawnPromotion { from: from_coordinate, to: coordinate_above, promotion_to: PieceType::Knight });
//            }
//
//            return moves
//        }
//
//        if self.board[coordinate_above as usize].is_empty() {
//            moves.push(ChessMove::Standard { from: from_coordinate, to: coordinate_above });
//
//            if from_coordinate.get_rank() == Rank::Second {
//                let coordinate_two_above = (coordinate_above + Direction::Up).unwrap();
//                if self.board[coordinate_two_above as usize].is_empty() {
//                    moves.push(ChessMove::Standard { from: from_coordinate, to: coordinate_two_above });
//                };
//            };
//        };
//
//
//
//        moves
//    }

//    fn get_possible_king_moves(&self, from_coordinate: Coordinate) -> Vec<ChessMove> {
//        Vec::new()
//    }
//
//    fn current_player_can_castle_kingside(&self) -> bool {
//        if !self.castling_rights[self.active_color].kingside {
//            return false
//        };
//
//        let key_squares = match self.active_color {
//            Color::White => [self.board[61], self.board[62]],
//            Color::Black => [self.board[5], self.board[6]]
//        };
//
//        key_squares.iter().fold(
//            true,
//            |acc, square| acc && square.is_empty() && !square.is_attacked_by(self.active_color.get_opposite())
//        )
//    }
//
//    fn current_player_can_castle_queenside(&self) -> bool {
//        if !self.castling_rights[self.active_color].queenside {
//            return false
//        };
//
//        let key_squares = match self.active_color {
//            Color::White => [self.board[57], self.board[58], self.board[59]],
//            Color::Black => [self.board[1], self.board[2], self.board[3]]
//        };
//
//        key_squares.iter().fold(
//            true,
//            |acc, square| acc && square.is_empty() && !square.is_attacked_by(self.active_color.get_opposite())
//        )
//    }
//
//    fn get_possible_queen_moves(&self, from_coordinate: Coordinate) -> Vec<ChessMove> {
//        Vec::new()
//    }
//
//    fn get_possible_bishop_moves(&self, from_coordinate: Coordinate) -> Vec<ChessMove> {
//        Vec::new()
//    }
//
//    fn get_possible_knight_moves(&self, from_coordinate: Coordinate) -> Vec<ChessMove> {
//        Vec::new()
//    }
//
//    fn get_possible_rook_moves(&self, from_coordinate: Coordinate) -> Vec<ChessMove> {
//        Vec::new()
//    }

    pub fn print(&self) {
//        let mut col_counter = 0;
//
//        print!("+---+---+---+---+---+---+---+---+\n");
//
//        for square in self.board.iter() {
//            print!("| {} ", square);
//
//            col_counter += 1;
//            col_counter = col_counter % 8;
//
//            if col_counter == 0 {
//                print!("|\n+---+---+---+---+---+---+---+---+\n");
//            }
//        }
//        io::stdout().flush().unwrap();
//
//        println!("{} to move.", self.active_color);
//        println!("White castling rights: {}.", self.castling_rights[Color::White]);
//        println!("Black castling rights: {}.", self.castling_rights[Color::Black]);
    }
}
