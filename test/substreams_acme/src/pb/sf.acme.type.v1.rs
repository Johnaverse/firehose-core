// @generated
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BlockHeader {
    #[prost(uint64, tag="1")]
    pub height: u64,
    #[prost(string, tag="2")]
    pub hash: ::prost::alloc::string::String,
    #[prost(uint64, optional, tag="3")]
    pub previous_num: ::core::option::Option<u64>,
    #[prost(string, optional, tag="4")]
    pub previous_hash: ::core::option::Option<::prost::alloc::string::String>,
    #[prost(uint64, tag="5")]
    pub final_num: u64,
    #[prost(string, tag="6")]
    pub final_hash: ::prost::alloc::string::String,
    #[prost(uint64, tag="7")]
    pub timestamp: u64,
}
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Block {
    #[prost(message, optional, tag="1")]
    pub header: ::core::option::Option<BlockHeader>,
    #[prost(message, repeated, tag="2")]
    pub transactions: ::prost::alloc::vec::Vec<Transaction>,
}
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Transaction {
    #[prost(string, tag="1")]
    pub r#type: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub hash: ::prost::alloc::string::String,
    #[prost(string, tag="3")]
    pub sender: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub receiver: ::prost::alloc::string::String,
    #[prost(message, optional, tag="5")]
    pub amount: ::core::option::Option<BigInt>,
    #[prost(message, optional, tag="6")]
    pub fee: ::core::option::Option<BigInt>,
    #[prost(bool, tag="7")]
    pub success: bool,
    #[prost(message, repeated, tag="8")]
    pub events: ::prost::alloc::vec::Vec<Event>,
}
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Event {
    #[prost(string, tag="1")]
    pub r#type: ::prost::alloc::string::String,
    #[prost(message, repeated, tag="2")]
    pub attributes: ::prost::alloc::vec::Vec<Attribute>,
}
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Attribute {
    #[prost(string, tag="1")]
    pub key: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub value: ::prost::alloc::string::String,
}
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BigInt {
    #[prost(bytes="vec", tag="1")]
    pub bytes: ::prost::alloc::vec::Vec<u8>,
}
// @@protoc_insertion_point(module)
