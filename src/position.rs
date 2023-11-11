use crate::board::Square;
use crate::chess_move::ChessMove;
use crate::color::Color;
use crate::castling_rights::CastlingRights;
use std::io::{self, Write};


pub struct Position {
    pub board: [Square; 64],
    pub active_color: Color,
    pub castling_rights: CastlingRights
}

impl Position {
    pub fn get_possible_moves(&self) -> Vec<ChessMove> {
        let mut moves = Vec::new();

        if self.current_player_can_castle_kingside() {
            moves.push(ChessMove::CastleKingside)
        };

        if self.current_player_can_castle_queenside() {
            moves.push(ChessMove::CastleQueenside)
        };

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
            |acc, square| acc && square.is_empty() && square.is_attacked_by(self.active_color.get_opposite())
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
            |acc, square| acc && square.is_empty() && square.is_attacked_by(self.active_color.get_opposite())
        )
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
