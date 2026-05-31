import Foundation

public struct User {
    public let id: String
    public var name: String
    public var email: String
}

public struct Product {
    public let id: String
    public var name: String
    public var price: Double
}

public enum Status {
    case active
    case inactive
}

public protocol Repository {
    func findById(_ id: String) -> User?
    func save(_ user: User) throws
}
