import Foundation

public class UserService {
    private var users: [String: User] = [:]

    public init() {}

    public func createUser(name: String, email: String) -> User {
        let id = generateId()
        let user = User(id: id, name: name, email: email)
        users[id] = user
        return user
    }

    public func getUser(id: String) -> User? {
        return users[id]
    }

    private func generateId() -> String {
        return UUID().uuidString
    }
}
