"use client"

import { useState } from "react"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardFooter } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { toast } from "sonner"

export function LoginForm() {
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [loading, setLoading] = useState(false)

    const handleLogin = async (backend: string) => {
        setLoading(true)
        const startTime = performance.now()
        try {
            let apiUrl

            // Select the appropriate backend API URL
            switch (backend) {
                case "axum":
                    apiUrl = `${process.env.NEXT_PUBLIC_AXUM_API_URL}/login`
                    break
                case "rocket":
                    apiUrl = `${process.env.NEXT_PUBLIC_ROCKET_API_URL}/login`
                    break
                case "gofiber":
                    apiUrl = `${process.env.NEXT_PUBLIC_GOFIBER_API_URL}/login`
                    break
                default:
                    throw new Error("Invalid backend")
            }

            const response = await fetch(apiUrl, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ username, password }),
            })

            const data = await response.json()
            const endTime = performance.now()
            const loginTimeMs = endTime - startTime
            const loginTime = `${Math.floor(loginTimeMs / 1000)}.${Math.floor(loginTimeMs % 1000)}` // seconds.milliseconds

            if (!data.success) {
                throw new Error(data.message || "Login failed")
            }

            // Store the token in localStorage
            localStorage.setItem(`${backend}_token`, data.token)

            toast("Login successful", {
                description: `Logged in using ${backend} backend in ${loginTime} seconds`,
            })

            // Remove the redirect to dashboard
            // router.push(`/dashboard?backend=${backend}`)
        } catch (error) {
            toast("Login failed", {
                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                // @ts-expect-error
                description: error.message,
            })
        } finally {
            setLoading(false)
        }
    }

    return (
        <Card>
            <Tabs defaultValue="axum" className={"p-5"}>
                <TabsList className="grid w-full grid-cols-3">
                    <TabsTrigger value="axum">Axum</TabsTrigger>
                    <TabsTrigger value="rocket">Rocket</TabsTrigger>
                    <TabsTrigger value="gofiber">GoFiber</TabsTrigger>
                </TabsList>
                <TabsContent value="axum">
                    <CardContent className="space-y-4 pt-4">
                        <div className="space-y-2">
                            <Label htmlFor="username-axum">Username</Label>
                            <Input
                                id="username-axum"
                                value={username}
                                onChange={(e) => setUsername(e.target.value)}
                                placeholder="Enter your username"
                            />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="password-axum">Password</Label>
                            <Input
                                id="password-axum"
                                type="password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                placeholder="Enter your password"
                            />
                        </div>
                    </CardContent>
                    <CardFooter className={"mt-5"}>
                        <Button className="w-full" onClick={() => handleLogin("axum")} disabled={loading}>
                            {loading ? "Logging in..." : "Login with Axum"}
                        </Button>
                    </CardFooter>
                </TabsContent>
                <TabsContent value="rocket">
                    <CardContent className="space-y-4 pt-4">
                        <div className="space-y-2">
                            <Label htmlFor="username-rocket">Username</Label>
                            <Input
                                id="username-rocket"
                                value={username}
                                onChange={(e) => setUsername(e.target.value)}
                                placeholder="Enter your username"
                            />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="password-rocket">Password</Label>
                            <Input
                                id="password-rocket"
                                type="password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                placeholder="Enter your password"
                            />
                        </div>
                    </CardContent>
                    <CardFooter className={"mt-5"}>
                        <Button className="w-full" onClick={() => handleLogin("rocket")} disabled={loading}>
                            {loading ? "Logging in..." : "Login with Rocket"}
                        </Button>
                    </CardFooter>
                </TabsContent>
                <TabsContent value="gofiber">
                    <CardContent className="space-y-4 pt-4">
                        <div className="space-y-2">
                            <Label htmlFor="username-gofiber">Username</Label>
                            <Input
                                id="username-gofiber"
                                value={username}
                                onChange={(e) => setUsername(e.target.value)}
                                placeholder="Enter your username"
                            />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="password-gofiber">Password</Label>
                            <Input
                                id="password-gofiber"
                                type="password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                placeholder="Enter your password"
                            />
                        </div>
                    </CardContent>
                    <CardFooter className={"mt-5"}>
                        <Button className="w-full" onClick={() => handleLogin("gofiber")} disabled={loading}>
                            {loading ? "Logging in..." : "Login with GoFiber"}
                        </Button>
                    </CardFooter>
                </TabsContent>
            </Tabs>
        </Card>
    )
}
