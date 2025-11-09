package utils

// 默认管理员账号配置（开发环境使用）
const (
	DEFAULT_ADMIN_USERNAME = "admin"
	DEFAULT_ADMIN_PASSWORD = "admin123"
)

// GetAdminCredentials 获取管理员账号配置
// 优先从环境变量获取，回退到默认值
func GetAdminCredentials() (string, string) {
	username := GetEnv("ADMIN_USERNAME")
	password := GetEnv("ADMIN_PASSWORD")

	// 如果环境变量未设置，使用默认值
	if username == "" {
		username = DEFAULT_ADMIN_USERNAME
	}
	if password == "" {
		password = DEFAULT_ADMIN_PASSWORD
	}

	return username, password
}

// ValidateAdminCredentials 验证管理员凭据
func ValidateAdminCredentials(inputUsername, inputPassword string) bool {
	// 检查输入是否为空
	if inputUsername == "" || inputPassword == "" {
		return false
	}

	// 获取配置的管理员账号
	adminUsername, adminPassword := GetAdminCredentials()

	// 验证用户名和密码
	return inputUsername == adminUsername && inputPassword == adminPassword
}

// GetAdminUsername 获取当前配置的主管理员用户名（用于JWT claims）
func GetAdminUsername(inputUsername string) string {
	// 如果验证通过，返回输入的用户名
	// 这样支持多管理员时每个人都有自己的身份
	return inputUsername
}

// GenerateAdminID 为管理员生成一个固定ID（用于JWT claims）
func GenerateAdminID(username string) string {
	// 基于用户名生成一个固定的ID，确保同一用户名总是得到相同的ID
	// 这里使用简单的哈希方式，也可以用更复杂的算法
	return "admin-" + username
}