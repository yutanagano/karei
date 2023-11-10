use std::fmt;


#[derive(Copy, Clone, Debug, PartialEq)]
pub enum Color {
    White,
    Black
}

impl fmt::Display for Color {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let as_str = match self {
            Color::White => "white",
            Color::Black => "black"
        };
        write!(f, "{as_str}")
    }
}
