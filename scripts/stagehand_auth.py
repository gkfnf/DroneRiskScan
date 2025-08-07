#!/usr/bin/env python3
"""
Stagehand Browser Automation Script for DroneRiskScan
使用真正的Stagehand SDK进行AI驱动的浏览器自动化认证
"""

import sys
import json
import asyncio
import os
from typing import Dict, Any, List, Optional
from datetime import datetime

try:
    # 使用Playwright直接实现"Stagehand风格"的AI驱动浏览器自动化
    from playwright.async_api import async_playwright, Page, Browser, BrowserContext
    import logging
    from urllib.parse import urlparse, urljoin
except ImportError as e:
    print(f"ERROR: Missing required packages: {e}", file=sys.stderr)
    print("Please install: pip install playwright", file=sys.stderr)
    sys.exit(1)

# 配置日志
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger('stagehand_auth')


class StagehandAuthenticator:
    """DroneRiskScan AI驱动的浏览器自动化认证器 - 使用Playwright实现类似Stagehand的功能"""
    
    def __init__(self, config: Dict[str, Any]):
        self.config = config
        self.playwright = None
        self.browser = None
        self.context = None
        self.page = None
        self.cookies = []
        self.session_data = {}
        
    async def initialize(self):
        """初始化浏览器自动化"""
        try:
            # 启动Playwright
            logger.info("Initializing Playwright browser automation...")
            self.playwright = await async_playwright().start()
            
            # 启动浏览器
            self.browser = await self.playwright.chromium.launch(
                headless=self.config.get("headless", True),
                args=['--disable-web-security', '--disable-features=VizDisplayCompositor']
            )
            
            # 创建浏览器上下文
            self.context = await self.browser.new_context(
                viewport={'width': 1920, 'height': 1080},
                user_agent='DroneRiskScan/1.0 AI-Browser-Automation'
            )
            
            # 创建新页面
            self.page = await self.context.new_page()
            
            logger.info("Browser automation initialized successfully")
            
        except Exception as e:
            logger.error(f"Failed to initialize browser automation: {e}")
            raise
    
    async def authenticate_bwapp(self, target_url: str, username: str, password: str) -> Dict[str, Any]:
        """使用智能浏览器自动化进行bWAPP认证"""
        try:
            result = {
                "success": False,
                "cookies": [],
                "session_data": {},
                "error": None,
                "screenshots": [],
                "network_logs": []
            }
            
            # 解析目标URL，确定登录页面
            parsed_url = urlparse(target_url)
            login_url = f"{parsed_url.scheme}://{parsed_url.netloc}/login.php"
            
            logger.info(f"Navigating to bWAPP login page: {login_url}")
            
            # 1. 导航到登录页面
            await self.page.goto(login_url, wait_until="networkidle")
            logger.info("Successfully navigated to login page")
            
            # 等待页面完全加载
            await self.page.wait_for_load_state('networkidle')
            
            # 2. 智能填写登录表单
            logger.info("Intelligently filling login form...")
            
            # 查找并填写用户名字段
            try:
                # 尝试多种可能的用户名字段选择器
                username_selectors = [
                    'input[name="login"]',
                    'input[name="username"]', 
                    'input[name="user"]',
                    'input[type="text"]',
                    '#login',
                    '#username'
                ]
                
                username_filled = False
                for selector in username_selectors:
                    try:
                        if await self.page.query_selector(selector):
                            await self.page.fill(selector, username)
                            logger.info(f"Filled username using selector: {selector}")
                            username_filled = True
                            break
                    except:
                        continue
                
                if not username_filled:
                    raise Exception("Could not find username field")
                
            except Exception as e:
                logger.error(f"Failed to fill username: {e}")
                result["error"] = f"Failed to fill username: {e}"
                return result
            
            # 查找并填写密码字段
            try:
                password_selectors = [
                    'input[name="password"]',
                    'input[type="password"]',
                    '#password'
                ]
                
                password_filled = False
                for selector in password_selectors:
                    try:
                        if await self.page.query_selector(selector):
                            await self.page.fill(selector, password)
                            logger.info(f"Filled password using selector: {selector}")
                            password_filled = True
                            break
                    except:
                        continue
                
                if not password_filled:
                    raise Exception("Could not find password field")
                
            except Exception as e:
                logger.error(f"Failed to fill password: {e}")
                result["error"] = f"Failed to fill password: {e}"
                return result
            
            # 设置安全级别为低（如果存在）
            try:
                security_select = await self.page.query_selector('select[name="security_level"]')
                if security_select:
                    await self.page.select_option('select[name="security_level"]', '0')
                    logger.info("Set security level to 'low' (0)")
                else:
                    logger.info("Security level selector not found, skipping")
            except Exception as e:
                logger.warning(f"Could not set security level: {e}")
            
            # 3. 点击登录按钮
            try:
                login_button_selectors = [
                    'button[type="submit"]',
                    'input[type="submit"]',
                    'button[name="form"]',
                    'input[name="form"]',
                    '.login-button',
                    'button:has-text("Login")',
                    'input[value="Login"]'
                ]
                
                button_clicked = False
                for selector in login_button_selectors:
                    try:
                        button = await self.page.query_selector(selector)
                        if button:
                            await button.click()
                            logger.info(f"Clicked login button using selector: {selector}")
                            button_clicked = True
                            break
                    except:
                        continue
                
                if not button_clicked:
                    # 尝试提交表单
                    form = await self.page.query_selector('form')
                    if form:
                        await form.evaluate('form => form.submit()')
                        logger.info("Submitted form directly")
                    else:
                        raise Exception("Could not find login button or form")
                
            except Exception as e:
                logger.error(f"Failed to submit login form: {e}")
                result["error"] = f"Failed to submit login form: {e}"
                return result
            
            # 等待登录处理
            await asyncio.sleep(3)
            await self.page.wait_for_load_state('networkidle')
            
            # 4. 验证登录是否成功
            current_url = self.page.url
            logger.info(f"Current URL after login: {current_url}")
            
            # 获取页面内容进行分析
            page_content = await self.page.content()
            page_title = await self.page.title()
            
            # 智能判断登录是否成功
            success_conditions = [
                "Choose your bug" in page_content,
                "portal.php" in current_url,
                "Welcome" in page_content,
                "Logout" in page_content,
                "dashboard" in current_url.lower(),
                "menu" in page_content.lower(),
                'name="login"' not in page_content  # 不再显示登录表单
            ]
            
            failure_conditions = [
                "bWAPP - Login" in page_title,
                'name="login"' in page_content,
                "login.php" in current_url,
                "Invalid credentials" in page_content,
                "Login failed" in page_content
            ]
            
            success_count = sum(1 for condition in success_conditions if condition)
            failure_count = sum(1 for condition in failure_conditions if condition)
            
            logger.info(f"Success indicators: {success_count}, Failure indicators: {failure_count}")
            
            if success_count > failure_count and success_count > 0:
                result["success"] = True
                logger.info("✅ Login successful!")
                
                # 提取Cookies
                cookies = await self.context.cookies()
                result["cookies"] = []
                for cookie in cookies:
                    result["cookies"].append({
                        "name": cookie["name"],
                        "value": cookie["value"],
                        "domain": cookie.get("domain", ""),
                        "path": cookie.get("path", "/"),
                        "secure": cookie.get("secure", False),
                        "httpOnly": cookie.get("httpOnly", False)
                    })
                
                logger.info(f"Extracted {len(result['cookies'])} cookies")
                for cookie in result["cookies"]:
                    logger.info(f"Cookie: {cookie['name']}={cookie['value']} (Domain: {cookie['domain']})")
                
                # 提取会话数据
                result["session_data"] = {
                    "current_url": current_url,
                    "page_title": page_title,
                    "success_indicators": success_count,
                    "failure_indicators": failure_count,
                    "timestamp": datetime.now().isoformat()
                }
                
            else:
                result["error"] = f"Login failed - success_indicators: {success_count}, failure_indicators: {failure_count}"
                logger.warning(f"❌ Login failed: {result['error']}")
                
                # 保存页面截图用于调试
                try:
                    screenshot_path = f"/tmp/bwapp_login_failed_{datetime.now().strftime('%Y%m%d_%H%M%S')}.png"
                    await self.page.screenshot(path=screenshot_path)
                    logger.info(f"Saved failure screenshot: {screenshot_path}")
                    result["screenshots"] = [screenshot_path]
                except Exception as screenshot_error:
                    logger.warning(f"Could not save screenshot: {screenshot_error}")
            
            return result
            
        except Exception as e:
            logger.error(f"Authentication error: {e}")
            return {
                "success": False,
                "error": str(e),
                "cookies": [],
                "session_data": {}
            }
    
    async def authenticate_generic(self, target_url: str, username: str, password: str) -> Dict[str, Any]:
        """通用Web应用认证 - 智能表单识别和填写"""
        try:
            logger.info(f"Starting generic authentication for: {target_url}")
            
            result = {
                "success": False,
                "cookies": [],
                "session_data": {},
                "error": None
            }
            
            # 导航到目标URL
            await self.page.goto(target_url, wait_until="networkidle")
            await self.page.wait_for_load_state('networkidle')
            
            # 智能查找并填写表单
            logger.info("Looking for login form...")
            
            # 查找用户名字段 - 更广泛的选择器
            username_selectors = [
                'input[name*="user"]',
                'input[name*="login"]', 
                'input[name*="email"]',
                'input[type="text"]:not([name*="search"]):not([name*="query"])',
                'input[type="email"]',
                '#username',
                '#login',
                '#email'
            ]
            
            username_filled = False
            for selector in username_selectors:
                try:
                    elements = await self.page.query_selector_all(selector)
                    if elements:
                        await elements[0].fill(username)
                        logger.info(f"Filled username using selector: {selector}")
                        username_filled = True
                        break
                except:
                    continue
            
            # 查找密码字段
            password_selectors = [
                'input[type="password"]',
                'input[name*="password"]',
                'input[name*="pass"]',
                '#password'
            ]
            
            password_filled = False
            for selector in password_selectors:
                try:
                    elements = await self.page.query_selector_all(selector)
                    if elements:
                        await elements[0].fill(password)
                        logger.info(f"Filled password using selector: {selector}")
                        password_filled = True
                        break
                except:
                    continue
            
            if not username_filled or not password_filled:
                result["error"] = "Could not find login form fields"
                return result
            
            # 提交表单
            submit_selectors = [
                'button[type="submit"]',
                'input[type="submit"]',
                'button:has-text("Login")',
                'button:has-text("Sign In")',
                'button:has-text("Submit")',
                '.login-btn',
                '.submit-btn'
            ]
            
            submitted = False
            for selector in submit_selectors:
                try:
                    button = await self.page.query_selector(selector)
                    if button:
                        await button.click()
                        logger.info(f"Clicked submit button using selector: {selector}")
                        submitted = True
                        break
                except:
                    continue
            
            if not submitted:
                # 尝试直接提交表单
                try:
                    await self.page.keyboard.press('Enter')
                    logger.info("Submitted form using Enter key")
                except:
                    result["error"] = "Could not submit login form"
                    return result
            
            # 等待处理
            await asyncio.sleep(3)
            await self.page.wait_for_load_state('networkidle')
            
            # 分析登录结果
            current_url = self.page.url
            page_content = await self.page.content()
            page_title = await self.page.title()
            
            # 简单的成功/失败判断
            success_patterns = ["welcome", "dashboard", "logout", "profile", "account", "home"]
            failure_patterns = ["login", "sign in", "error", "invalid", "failed"]
            
            success_score = sum(1 for pattern in success_patterns if pattern.lower() in page_content.lower())
            failure_score = sum(1 for pattern in failure_patterns if pattern.lower() in page_content.lower())
            
            # URL变化也是成功的指标
            if target_url != current_url and "login" not in current_url.lower():
                success_score += 2
            
            if success_score > failure_score:
                result["success"] = True
                
                # 提取cookies
                cookies = await self.context.cookies()
                result["cookies"] = [
                    {
                        "name": cookie["name"],
                        "value": cookie["value"],
                        "domain": cookie.get("domain", ""),
                        "path": cookie.get("path", "/")
                    }
                    for cookie in cookies
                ]
                
                result["session_data"] = {
                    "current_url": current_url,
                    "page_title": page_title,
                    "success_score": success_score,
                    "failure_score": failure_score,
                    "timestamp": datetime.now().isoformat()
                }
                
                logger.info(f"✅ Generic login successful (score: {success_score} vs {failure_score})")
            else:
                result["error"] = f"Login failed (score: {success_score} vs {failure_score})"
                logger.warning(f"❌ {result['error']}")
            
            return result
            
        except Exception as e:
            logger.error(f"Generic authentication error: {e}")
            return {
                "success": False,
                "error": str(e),
                "cookies": [],
                "session_data": {}
            }
    
    async def discover_function_points(self, target_url: str) -> List[Dict[str, Any]]:
        """发现页面功能点 - 用于后续安全测试"""
        try:
            logger.info(f"Discovering function points for: {target_url}")
            
            await self.page.goto(target_url)
            await asyncio.sleep(2)
            
            # 使用AI分析页面功能点
            analysis = await self.page.extract({
                "instruction": """Analyze this web page and identify all potential security testing points including:
                1. Forms with input fields (especially search, contact, comment forms)
                2. Links with GET parameters 
                3. File upload features
                4. Database interaction points (search, filters, etc.)
                5. Administrative or privileged areas
                
                For each point, identify the URL, method, parameters, and potential vulnerability types.""",
                "schema": {
                    "type": "object",
                    "properties": {
                        "function_points": {
                            "type": "array",
                            "items": {
                                "type": "object",
                                "properties": {
                                    "url": {"type": "string"},
                                    "method": {"type": "string"},
                                    "type": {"type": "string"},
                                    "parameters": {"type": "array", "items": {"type": "string"}},
                                    "vulnerability_types": {"type": "array", "items": {"type": "string"}},
                                    "description": {"type": "string"}
                                }
                            }
                        }
                    }
                }
            })
            
            function_points = analysis.get("function_points", [])
            logger.info(f"Discovered {len(function_points)} function points")
            
            return function_points
            
        except Exception as e:
            logger.error(f"Function point discovery error: {e}")
            return []
    
    async def cleanup(self):
        """清理资源"""
        try:
            if self.browser:
                await self.browser.close()
                logger.info("Browser closed successfully")
            if self.playwright:
                await self.playwright.stop()
                logger.info("Playwright stopped successfully")
        except Exception as e:
            logger.error(f"Cleanup error: {e}")


async def main():
    """主函数 - 处理命令行参数"""
    if len(sys.argv) < 2:
        print("Usage: python stagehand_auth.py <command> [args...]", file=sys.stderr)
        print("Commands:", file=sys.stderr)
        print("  auth <target_url> <username> <password> [auth_type]", file=sys.stderr)
        print("  discover <target_url>", file=sys.stderr)
        sys.exit(1)
    
    command = sys.argv[1]
    
    try:
        if command == "auth":
            if len(sys.argv) < 5:
                print("Usage: python stagehand_auth.py auth <target_url> <username> <password> [auth_type]", file=sys.stderr)
                sys.exit(1)
            
            target_url = sys.argv[2]
            username = sys.argv[3]
            password = sys.argv[4]
            auth_type = sys.argv[5] if len(sys.argv) > 5 else "bwapp"
            
            config = {
                "headless": True,  # 可以改为False进行调试
                "timeout": 30
            }
            
            authenticator = StagehandAuthenticator(config)
            
            try:
                await authenticator.initialize()
                
                if auth_type == "bwapp":
                    result = await authenticator.authenticate_bwapp(target_url, username, password)
                else:
                    result = await authenticator.authenticate_generic(target_url, username, password)
                
                # 输出JSON结果给Go程序
                print(json.dumps(result))
                
            finally:
                await authenticator.cleanup()
        
        elif command == "discover":
            if len(sys.argv) < 3:
                print("Usage: python stagehand_auth.py discover <target_url>", file=sys.stderr)
                sys.exit(1)
            
            target_url = sys.argv[2]
            
            config = {"headless": True}
            authenticator = StagehandAuthenticator(config)
            
            try:
                await authenticator.initialize()
                function_points = await authenticator.discover_function_points(target_url)
                print(json.dumps({"function_points": function_points}))
                
            finally:
                await authenticator.cleanup()
        
        else:
            print(f"Unknown command: {command}", file=sys.stderr)
            sys.exit(1)
    
    except Exception as e:
        error_result = {
            "success": False,
            "error": str(e),
            "cookies": [],
            "session_data": {}
        }
        print(json.dumps(error_result))
        sys.exit(1)


if __name__ == "__main__":
    asyncio.run(main())