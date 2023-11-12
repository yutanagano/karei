use crate::board::{Coordinate, Square};
use crate::piece::PieceType;
use crate::chess_move::ChessMove;
use crate::color::Color;
use crate::castling_rights::CastlingRights;
use std::io::{self, Write};
use std::ops::{Index, IndexMut};


pub struct Position {
    board: [Square; 64],
    active_color: Color,
    castling_rights: CastlingRights
}

impl Position {
    pub fn new(board: [Square; 64], active_color: Color, castling_rights: CastlingRights) -> Self {
        Position { board, active_color, castling_rights }
    }

    pub fn get_possible_moves(&self) -> Vec<ChessMove> {
        let mut moves = Vec::new();

        if self.current_player_can_castle_kingside() {
            moves.push(ChessMove::CastleKingside)
        };

        if self.current_player_can_castle_queenside() {
            moves.push(ChessMove::CastleQueenside)
        };
        
        for (index, square) in self.board.iter().enumerate() {
            if square.is_empty() {
                continue;
            };

            let piece = square.get_piece().unwrap();

            if piece.color == self.active_color.get_opposite() {
                continue;
            };

            let pseudo_legal_moves = match piece.piece_type {
                PieceType::Pawn => self.get_possible_pawn_moves(index),
                PieceType::King => self.get_possible_king_moves(index),
                PieceType::Queen => self.get_possible_queen_moves(index),
                PieceType::Bishop => self.get_possible_bishop_moves(index),
                PieceType::Knight => self.get_possible_knight_moves(index),
                PieceType::Rook => self.get_possible_rook_moves(index)
            };

            moves.extend(pseudo_legal_moves);
        }

        moves
    }

    fn current_player_can_castle_kingside(&self) -> bool {
        if !self.castling_rights[self.active_color].kingside {
            return false
        };

        let key_squares = match self.active_color {
            Color::White => [self.board[61], self.board[62]],
            Color::Black => [self.board[5], self.board[6]]
        };

        key_squares.iter().fold(
            true,
            |acc, square| acc && square.is_empty() && !square.is_attacked_by(self.active_color.get_opposite())
        )
    }

    fn current_player_can_castle_queenside(&self) -> bool {
        if !self.castling_rights[self.active_color].queenside {
            return false
        };

        let key_squares = match self.active_color {
            Color::White => [self.board[57], self.board[58], self.board[59]],
            Color::Black => [self.board[1], self.board[2], self.board[3]]
        };

        key_squares.iter().fold(
            true,
            |acc, square| acc && square.is_empty() && !square.is_attacked_by(self.active_color.get_opposite())
        )
    }

    fn get_possible_pawn_moves(&self, from_square_index: usize) -> Vec<ChessMove> {
        Vec::new()
    }

    fn get_possible_king_moves(&self, from_square_index: usize) -> Vec<ChessMove> {
        Vec::new()
    }

    fn get_possible_queen_moves(&self, from_square_index: usize) -> Vec<ChessMove> {
        Vec::new()
    }

    fn get_possible_bishop_moves(&self, from_square_index: usize) -> Vec<ChessMove> {
        Vec::new()
    }

    fn get_possible_knight_moves(&self, from_square_index: usize) -> Vec<ChessMove> {
        Vec::new()
    }

    fn get_possible_rook_moves(&self, from_square_index: usize) -> Vec<ChessMove> {
        Vec::new()
    }

    pub fn print(&self) {
        let mut col_counter = 0;

        print!("+---+---+---+---+---+---+---+---+\n");

        for square in self.board.iter() {
            print!("| {} ", square);

            col_counter += 1;
            col_counter = col_counter % 8;

            if col_counter == 0 {
                print!("|\n+---+---+---+---+---+---+---+---+\n");
            }
        }
        io::stdout().flush().unwrap();

        println!("{} to move.", self.active_color);
        println!("White castling rights: {}.", self.castling_rights[Color::White]);
        println!("Black castling rights: {}.", self.castling_rights[Color::Black]);
    }
}

impl Index<Coordinate> for Position {
    type Output = Square;

    fn index(&self, index: Coordinate) -> &Self::Output {
        &self.board[index as usize]
    }
}

impl IndexMut<Coordinate> for Position {
    fn index_mut(&mut self, index: Coordinate) -> &mut Self::Output {
        &mut self.board[index as usize]
    }
}
