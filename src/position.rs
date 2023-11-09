use crate::board::Board;
use crate::color::Color;
use crate::castling_rights::CastlingRights;
use std::io::{self, Write};


pub struct Position {
    board: Board,
    active_color: Color,
    castling_rights: CastlingRights
}

impl Position {
    pub fn new(board: Board, active_color: Color, castling_rights: CastlingRights) -> Self {
        Position {
            board,
            active_color,
            castling_rights
        }
    }

    pub fn print(&self) {
        let mut col_counter = 0;

        print!("+---+---+---+---+---+---+---+---+\n");

        for square in self.board.iter_squares() {
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
