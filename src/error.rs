use thiserror::Error;

pub type Result<T> = std::result::Result<T, Error>;

#[derive(Error,Debug)]
pub enum Error{
    
    // Third party errors
    #[error("IO error: {0}")]
    Io(#[from] std::io::Error),
    #[error("Var error: {0}")]
    Var(#[from] std::env::VarError),
    #[error("Reqwest error: {0}")]
    Reqwest(#[from] reqwest::Error),
}