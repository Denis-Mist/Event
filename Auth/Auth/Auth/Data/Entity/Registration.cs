using System.ComponentModel.DataAnnotations;
using Microsoft.VisualBasic;

namespace Auth.Data.Entity
{
    public class Registration
    {
        public int Id { get; set; }
        [Required] public string Name { get; set; } = string.Empty;

        [Required] public string Phone { get; set; } = string.Empty;

        [Required] public string Email { get; set; } = string.Empty;

        [Required] public string Password { get; set; } = string.Empty;

    }
}