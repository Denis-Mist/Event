namespace Auth.Data.Entity
{
    public class UserRole
    {
        public int Id { get; set; }
        public string? RoleName { get; set; }
        public int UserId { get; set; }
    }
}