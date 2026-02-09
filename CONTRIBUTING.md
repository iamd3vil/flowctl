# Contributing

Thank you for your interest in contributing to flowctl!

## Getting Started

### 1. Check Existing Issues

Before starting work, check the [issue tracker](https://github.com/cvhariharan/flowctl/issues) to see if someone is already working on a similar feature or bug fix.

### 2. Create or Comment on an Issue

- For bugs: Create a new issue describing the problem, steps to reproduce, and expected behavior
- For features: Create a new issue explaining the feature, use case, and proposed implementation
- For existing issues: Comment on the issue expressing your interest in working on it

### 3. Wait for Assignment

Pull requests will only be accepted for assigned issues.

## Development

### Discuss Before Implementing

Before submitting a PR, discuss the implementation details in the issue. This includes:
- Overall approach and architecture
- Which files/modules will be affected
- Any breaking changes or compatibility concerns

This discussion helps ensure your work aligns with the project's goals and saves time for both contributors and maintainers.

### Setting Up Your Development Environment

1. Fork the repository
2. Clone your fork:
```bash
   git clone https://github.com/YOUR_USERNAME/flowctl.git
   cd flowctl
```
3. Refer to [development setup](https://flowctl.net/docs/development/setup/)
```bash
   make dev-docker
```
4. Create a new branch for your work:
```bash
   git checkout -b feature/your-feature-name
```

### Code Quality Standards

- Write meaningful commit messages
- Follow best practices
- Update documentation if your changes affect user-facing features
- Keep changes focused: one PR should address one issue

### Pull Request Guidelines

#### Size and Scope

- Keep PRs focused and reasonably sized
- If your work requires many file changes, break it into smaller, logical PRs

**PRs that show signs of being mass-generated without careful review will be closed.** This includes:
- PRs with excessive file changes that lack clear justification
- Changes that don't address the specific issue
- Code that hasn't been tested or understood by the contributor
- Generic refactoring across many files without prior discussion

#### Submitting Your PR

1. Push your branch to your fork
2. Open a pull request against the `master` branch
3. Reference the issue number in your PR description (e.g., "Fixes #123")
4. Provide a clear description of:
   - What the PR does
   - How it addresses the issue
   - Any testing you've performed
   - Screenshots/demos for UI changes

## Code of Conduct

- Be respectful and constructive in all interactions
- Help create a welcoming environment for all contributors
- Focus on the code and ideas, not individuals

## Getting Help

If you have questions:
- Ask in the relevant issue thread
- Check the [documentation](https://flowctl.net/docs)
- Use [discussions](https://github.com/cvhariharan/flowctl/discussions)
