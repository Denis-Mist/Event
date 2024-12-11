using System.ComponentModel.DataAnnotations;

namespace Auth.Data.Entity
{
    public class Login
    {
        [Required]
        [EmailAddress]
        public string? Email { get; set; }
        [Required]
        public string? Password { get; set; }
    }
}

