# trinity-proto
Proto for the Trinity project

# TODO
- [ ] Session-based auth
- [ ] ABAC authorization
- [ ] golint and build + tests CI

# Use Cases
- User
    - [ ] Get User by ID
    - [ ] Create User
    - [ ] Promote User
    - [ ] Demote User
    - [ ] Delete User

- Auth
    - [ ] Login
    - [ ] Logout
    - [ ] Verify
    - [ ] Logout specific session

# Talking with the outside 
In terms of visibility, a module is allowed to import and use other modules' clients. And that's the only single piece of code they can import from the other modules. Ideally a client is defined as an interface, allowing to go with a direct code call implementation or an over-the-network implementation, in case it's needed (for instance, by an actual external application). 