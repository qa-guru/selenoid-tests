package hubapi

// CreateSessionBody builds a W3C New Session body (HubSessionApi.createSessionBody).
func CreateSessionBody(browserName, browserVersion string) map[string]any {
	return map[string]any{
		"capabilities": map[string]any{
			"alwaysMatch": createAlwaysMatch(browserName, browserVersion),
		},
	}
}

func createAlwaysMatch(browserName, browserVersion string) map[string]any {
	alwaysMatch := map[string]any{
		"browserName":    browserName,
		"browserVersion": browserVersion,
	}
	switch browserName {
	case "firefox":
		alwaysMatch["moz:firefoxOptions"] = map[string]any{"args": dockerFirefoxArgs()}
	case "msedge":
		alwaysMatch["ms:edgeOptions"] = map[string]any{"args": dockerEdgeArgs()}
	default:
		alwaysMatch["goog:chromeOptions"] = map[string]any{"args": dockerChromeArgs()}
	}
	return alwaysMatch
}

func dockerChromeArgs() []string {
	return []string{
		"--headless=new",
		"--no-sandbox",
		"--disable-dev-shm-usage",
	}
}

func dockerFirefoxArgs() []string {
	return []string{"-headless"}
}

func dockerEdgeArgs() []string {
	return []string{
		"--headless=new",
		"--no-sandbox",
		"--disable-dev-shm-usage",
	}
}
