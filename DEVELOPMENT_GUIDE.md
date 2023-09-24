# Development Guidelines
Welcome to the development community of Diggity! These guidelines are designed to help you contribute effectively to our project and ensure a smooth and collaborative development process.

## Prerequisites
Before you begin, ensure that you have the following prerequisites installed on your system:

- Go (version 1.9 or higher)
- Git
- Your Preferred Text Editor or IDE

## Build and Run
- First, clone the Diggity repository to your local machine using Git:
```bash
git clone https://github.com/carbonetes/diggity.git
cd diggity
```
- Install any project-specific dependencies (if applicable).
```bash
go mod tidy
```
- Build the binary.
```bash
go build
```
- Run the binary.
```bash
./diggity
```
This will start the application and make it accessible locally.

## Version Control
- We use Git for version control. Familiarize yourself with Git if you aren't already.

## Coding Standards
- Follow the coding style and conventions for the project language (e.g., Go) as outlined in our project's specific style guide.

## Branching Strategy
- Create feature branches for your work and make pull requests (PRs) to the main branch.
- Branch names should be clear and descriptive, e.g., feature/new-feature or bugfix/issue-fix.

## Commit Guidelines
- Each commit should be atomic and focused on a single task.
- Use clear and concise commit messages with a descriptive summary.<br /> 
**Note**: Please follow the [Karma 6.4](http://karma-runner.github.io/6.4/dev/git-commit-msg.html) format.
- Reference related issues or PRs in your commits when applicable, e.g., "Fix #123: Updated authentication logic."

## Code Reviews
- All code changes must undergo a code review before being merged.
- Be open to feedback and constructive criticism during code reviews.
- Address review comments and make necessary changes before merging.

## Testing
- Write unit tests and, where applicable, integration tests for your code.
- Ensure that all tests pass before submitting a PR. To run the project's tests, use the following command:
```bash
go test ./...
```
- Add tests for bug fixes or new features, and update existing tests as needed.

## Documentation
- Keep code comments, function/method signatures, and user-facing documentation up to date.
- Document code that is not self-explanatory or may be unclear to others.
- Maintain a clear and organized project documentation.

## Continuous Integration (CI)
- Our CI pipeline runs automated tests, code linting, and other checks.
- Ensure that your code passes all CI checks before creating a PR.

## Reporting Issues
- Use our GitHub Issue Tracker to report bugs, suggest enhancements, or submit feature requests.
- Provide detailed information, steps to reproduce, and any relevant context when creating issues.

## Getting Help
If you have questions, need assistance, or want to discuss potential contributions, please don't hesitate to contact us at [eng@carbonetes.com](mailto:eng@carbonetes.com).

## Conclusion
By following these Development Guidelines, you'll contribute to the success of Diggity and help maintain a collaborative and efficient development process. Thank you for your dedication to our project!