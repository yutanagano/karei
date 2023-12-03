use crate::board::{Coordinate, Direction, Rank, Board, File};
use crate::piece::PieceType;
use crate::chess_move::ChessMove;
use crate::color::Color;
use crate::castling_rights::CastlingRights;


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

        let files = [File::A, File::B, File::C, File::D, File::E, File::F, File::G, File::H];
        let ranks = [Rank::First, Rank::Second, Rank::Third, Rank::Fourth, Rank::Fifth, Rank::Sixth, Rank::Seventh, Rank::Eighth];

        for file in files.iter() {
            for rank in ranks.iter() {
                let coordinate = Coordinate::new(*file, *rank);
                let square = self.board[coordinate];

                if square.is_empty() {
                    continue;
                };

                let piece = square.get_piece().unwrap();

                if piece.color == self.active_color.get_opposite() {
                    continue;
                };

                let pseudo_legal_moves = match piece.piece_type {
                    PieceType::Pawn => self.get_possible_pawn_moves(coordinate),
                    PieceType::King => self.get_possible_king_moves(coordinate),
                    PieceType::Queen => self.get_possible_queen_moves(coordinate),
                    PieceType::Bishop => self.get_possible_bishop_moves(coordinate),
                    PieceType::Knight => self.get_possible_knight_moves(coordinate),
                    PieceType::Rook => self.get_possible_rook_moves(coordinate)
                };

                moves.extend(pseudo_legal_moves);
            }
        }

        moves
    }

    fn get_possible_pawn_moves(&self, from_coordinate: Coordinate) -> Vec<ChessMove> {
        let mut moves = Vec::new();
        let coordinate_above = (from_coordinate.try_moving(Direction::up())).unwrap();
        let capture_coordinates = [
            from_coordinate.try_moving(Direction::upper_left()).unwrap(),
            from_coordinate.try_moving(Direction::upper_right()).unwrap()
        ];

        if from_coordinate.get_rank() == Rank::Seventh {
            for capture_coordinate in capture_coordinates.into_iter() {
                if self.board[capture_coordinate].has_piece_of_color(self.active_color.get_opposite()) {
                    moves.push(ChessMove::PawnPromotion { from: from_coordinate, to: capture_coordinate, promotion_to: PieceType::Queen });
                    moves.push(ChessMove::PawnPromotion { from: from_coordinate, to: capture_coordinate, promotion_to: PieceType::Rook });
                    moves.push(ChessMove::PawnPromotion { from: from_coordinate, to: capture_coordinate, promotion_to: PieceType::Bishop });
                    moves.push(ChessMove::PawnPromotion { from: from_coordinate, to: capture_coordinate, promotion_to: PieceType::Knight });
                };
            };

            if self.board[coordinate_above].is_empty() {
                moves.push(ChessMove::PawnPromotion { from: from_coordinate, to: coordinate_above, promotion_to: PieceType::Queen });
                moves.push(ChessMove::PawnPromotion { from: from_coordinate, to: coordinate_above, promotion_to: PieceType::Rook });
                moves.push(ChessMove::PawnPromotion { from: from_coordinate, to: coordinate_above, promotion_to: PieceType::Bishop });
                moves.push(ChessMove::PawnPromotion { from: from_coordinate, to: coordinate_above, promotion_to: PieceType::Knight });
            };

            return moves
        };

        if self.board[coordinate_above].is_empty() {
            moves.push(ChessMove::Standard { from: from_coordinate, to: coordinate_above });

            if from_coordinate.get_rank() == Rank::Second {
                let coordinate_two_above = (coordinate_above.try_moving(Direction::up())).unwrap();
                if self.board[coordinate_two_above].is_empty() {
                    moves.push(ChessMove::Standard { from: from_coordinate, to: coordinate_two_above });
                };
            };
        };

        for capture_coordinate in capture_coordinates.into_iter() {
            if self.board[capture_coordinate].has_piece_of_color(self.active_color.get_opposite()) {
                moves.push(ChessMove::Standard { from: from_coordinate, to: capture_coordinate })
            };
        };

        moves
    }

    fn get_possible_king_moves(&self, from_coordinate: Coordinate) -> Vec<ChessMove> {
        Vec::new()
    }

    fn current_player_can_castle_kingside(&self) -> bool {
        if !self.castling_rights[self.active_color].kingside {
            return false
        };

        let key_squares = match self.active_color {
            Color::White => [self.board[Coordinate::new(File::F, Rank::First)], self.board[Coordinate::new(File::G, Rank::First)]],
            Color::Black => [self.board[Coordinate::new(File::F, Rank::Eighth)], self.board[Coordinate::new(File::G, Rank::Eighth)]]
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
            Color::White => [self.board[Coordinate::new(File::D, Rank::First)], self.board[Coordinate::new(File::C, Rank::First)], self.board[Coordinate::new(File::B, Rank::First)]],
            Color::Black => [self.board[Coordinate::new(File::D, Rank::Eighth)], self.board[Coordinate::new(File::C, Rank::Eighth)], self.board[Coordinate::new(File::B, Rank::Eighth)]]
        };

        key_squares.iter().fold(
            true,
            |acc, square| acc && square.is_empty() && !square.is_attacked_by(self.active_color.get_opposite())
        )
    }

    fn get_possible_queen_moves(&self, from_coordinate: Coordinate) -> Vec<ChessMove> {
        Vec::new()
    }

    fn get_possible_bishop_moves(&self, from_coordinate: Coordinate) -> Vec<ChessMove> {
        Vec::new()
    }

    fn get_possible_knight_moves(&self, from_coordinate: Coordinate) -> Vec<ChessMove> {
        Vec::new()
    }

    fn get_possible_rook_moves(&self, from_coordinate: Coordinate) -> Vec<ChessMove> {
        Vec::new()
    }

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
