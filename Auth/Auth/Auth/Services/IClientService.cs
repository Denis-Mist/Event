using Auth.Data.Entity;

namespace Auth.Services
{
    public interface IClientService
    {
        Task<(bool flag, string Message)> RegisterUserAsync(Registration model);
        Task<(bool flag, string Token)> LoginUserAsync(Login model);
    }
}