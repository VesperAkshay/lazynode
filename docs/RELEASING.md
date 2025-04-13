# Release Process

This document outlines the process for creating new releases of LazyNode.

## Prerequisites

- Access to the GitHub repository with permission to create releases
- Git installed locally
- Go 1.19+ installed

## Steps to Create a New Release

1. **Ensure the code is ready for release**
   - All tests pass: `make test`
   - The code builds successfully: `make build`
   - Documentation is up-to-date

2. **Update version information**
   - Update any version references in the codebase if necessary

3. **Create and push a git tag**
   ```bash
   # Tag the release (use semantic versioning)
   git tag -a v1.0.0 -m "Release v1.0.0"
   
   # Push the tag to GitHub
   git push origin v1.0.0
   ```

4. **Wait for GitHub Actions**
   - The GitHub Actions workflow will automatically:
     - Build binaries for all supported platforms
     - Create a GitHub release
     - Upload all artifacts to the release

5. **Verify the release**
   - Check the GitHub Actions logs for any issues
   - Download and test the released binaries
   - Ensure all artifacts were uploaded correctly

6. **Announce the release**
   - Update the website (if applicable)
   - Post release notes to relevant channels

## Manual Release Process

If you need to create a release manually:

1. **Build the binaries locally**
   ```bash
   # Clean any previous builds
   make clean
   
   # Build for all platforms
   make release
   ```

2. **Create a release on GitHub**
   - Go to the "Releases" section in the GitHub repository
   - Click "Draft a new release"
   - Select the tag version
   - Add release notes
   - Upload the binaries from the `dist` directory
   - Publish the release

## Troubleshooting

### Build Failures

If builds fail, check:
- Go version compatibility
- Dependencies and module setup
- Platform-specific issues in the code

### GitHub Actions Issues

If GitHub Actions fail:
- Check the workflow logs
- Ensure the repository has the correct permissions set
- Verify that the workflow file is correct

## Future Improvements

- Automate changelog generation
- Add automated testing on all platforms
- Set up continuous deployment to package managers 