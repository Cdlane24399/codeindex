import Foundation

func main() {
    let service = UserService()
    let user = service.createUser(name: "Alice", email: "alice@example.com")
    print("Created user: \(user.name)")
}

main()
