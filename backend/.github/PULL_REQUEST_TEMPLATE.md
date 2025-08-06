## Summary

Brief description of the changes in this PR.

## Type of Change

Please mark the relevant option:

- [ ] 🐛 Bug fix (non-breaking change which fixes an issue)
- [ ] ✨ New feature (non-breaking change which adds functionality)  
- [ ] 💥 Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] 🔧 Refactoring (code change that neither fixes a bug nor adds a feature)
- [ ] 📝 Documentation update
- [ ] 🎨 Style change (formatting, missing semi colons, etc; no code change)
- [ ] ⚡ Performance improvement
- [ ] 🧪 Test update
- [ ] 🏗️ Build/CI change

## Related Issues

Closes #(issue number)
Fixes #(issue number)
Related to #(issue number)

## Changes Made

### Domain Layer Changes
- [ ] Added/modified entities in `domain/{module}/entity/`
- [ ] Added/modified repository interfaces in `domain/{module}/repository/`
- [ ] Added/modified domain services in `domain/{module}/service/`
- [ ] Added domain-specific errors

Details:
- 

### Infrastructure Layer Changes
- [ ] Added/modified repository implementations in `infra/repository/`
- [ ] Added/modified database models/migrations
- [ ] Added/modified configuration
- [ ] Added/modified external service integrations

Details:
- 

### API Layer Changes
- [ ] Added/modified handlers in `api/handler/`
- [ ] Added/modified request/response models in `api/model/`
- [ ] Added/modified middleware in `api/middleware/`
- [ ] Added/modified routes in `api/router/`

Details:
- 

### Application Layer Changes
- [ ] Modified dependency injection in `application/application.go`
- [ ] Added/modified application services
- [ ] Added/modified use cases

Details:
- 

## Database Changes

- [ ] No database changes
- [ ] Schema changes (migrations included)
- [ ] Data migration required

If schema changes, please describe:

```sql
-- Example migration
CREATE TABLE example (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);
```

## API Changes

- [ ] No API changes
- [ ] New endpoints added
- [ ] Existing endpoints modified
- [ ] Breaking API changes

### New Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST   | `/api/v1/example` | Creates a new example |

### Modified Endpoints

| Method | Endpoint | Changes |
|--------|----------|---------|
| GET    | `/api/v1/example` | Added new query parameter |

### Breaking Changes

Describe any breaking changes and migration steps:

## Testing

- [ ] Unit tests added/updated
- [ ] Integration tests added/updated  
- [ ] API tests added/updated
- [ ] Manual testing completed
- [ ] All existing tests pass

### Test Coverage

- Overall coverage: `%`
- New code coverage: `%`

### Manual Testing

Describe the manual testing performed:

1. Test scenario 1:
   - Steps: 
   - Expected result:
   - Actual result:

## Security Considerations

- [ ] No security impact
- [ ] Security review completed
- [ ] Potential security implications (describe below)

If there are security implications, describe them:

## Performance Impact

- [ ] No performance impact
- [ ] Performance improvement
- [ ] Potential performance degradation (mitigated)
- [ ] Requires performance testing

If performance impact, describe:

## Documentation

- [ ] Code is self-documenting
- [ ] Inline comments added for complex logic
- [ ] README updated
- [ ] API documentation updated
- [ ] Architecture documentation updated

## Configuration Changes

- [ ] No configuration changes
- [ ] New configuration options added
- [ ] Existing configuration modified
- [ ] Breaking configuration changes

If configuration changes, document them:

```yaml
# New configuration
new_feature:
  enabled: true
  timeout: 30s
```

## Deployment Notes

- [ ] No special deployment requirements
- [ ] Requires environment variables
- [ ] Requires database migration
- [ ] Requires service restart
- [ ] Requires dependency updates

Special deployment steps:

## Checklist

### Before Submitting
- [ ] I have read the [contributing guidelines](CONTRIBUTING.md)
- [ ] My code follows the DDD architecture principles
- [ ] My code follows the project's coding style
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings or errors
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes

### Architecture Compliance
- [ ] Changes follow DDD layer boundaries
- [ ] Domain layer contains no infrastructure dependencies
- [ ] Repository interfaces are defined in domain layer
- [ ] Repository implementations are in infrastructure layer
- [ ] Business logic is in domain/application layers, not in handlers

### Code Quality
- [ ] Code is properly formatted (`go fmt`)
- [ ] Code passes linting (`golint`, `go vet`)
- [ ] No magic numbers or hardcoded values
- [ ] Error handling is comprehensive
- [ ] Logging is appropriate and structured

## Screenshots/Recordings

If applicable, add screenshots or recordings to help explain your changes:

## Additional Notes

Any additional information that reviewers should know:

---

**For Reviewers:**

Please review:
1. ✅ DDD architecture compliance
2. ✅ Code quality and style
3. ✅ Test coverage and quality
4. ✅ Documentation completeness
5. ✅ Security implications
6. ✅ Performance impact
7. ✅ Breaking changes documentation
