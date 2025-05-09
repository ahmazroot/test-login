"use client"

import type React from "react"

import { useState, useRef } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent,CardFooter } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { toast } from "sonner"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"

export function RegisterForm() {
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [profilePhoto, setProfilePhoto] = useState<File | null>(null)
    const [idPhoto, setIdPhoto] = useState<File | null>(null)
    const [loading, setLoading] = useState(false)
    const [profilePhotoPreview, setProfilePhotoPreview] = useState<string | null>(null)
    const [idPhotoPreview, setIdPhotoPreview] = useState<string | null>(null)

    const profilePhotoRef = useRef<HTMLInputElement>(null)
    const idPhotoRef = useRef<HTMLInputElement>(null)


    const handleProfilePhotoChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files[0]) {
            const file = e.target.files[0]
            setProfilePhoto(file)

            // Create preview
            const reader = new FileReader()
            reader.onload = (event) => {
                setProfilePhotoPreview(event.target?.result as string)
            }
            reader.readAsDataURL(file)
        }
    }

    const handleIdPhotoChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files[0]) {
            const file = e.target.files[0]
            setIdPhoto(file)

            // Create preview
            const reader = new FileReader()
            reader.onload = (event) => {
                setIdPhotoPreview(event.target?.result as string)
            }
            reader.readAsDataURL(file)
        }
    }

    const handleRegister = async (backend: string) => {
        if (!username || !password || !profilePhoto || !idPhoto) {
            toast("Registration failed", {
                description: "Please fill in all fields and upload both photos",
            })
            return
        }

        setLoading(true)
        const startTime = performance.now()

        try {
            let apiUrl

            // Select the appropriate backend API URL
            switch (backend) {
                case "axum":
                    apiUrl = `${process.env.NEXT_PUBLIC_AXUM_API_URL}/register`
                    break
                case "rocket":
                    apiUrl = `${process.env.NEXT_PUBLIC_ROCKET_API_URL}/register`
                    break
                case "gofiber":
                    apiUrl = `${process.env.NEXT_PUBLIC_GOFIBER_API_URL}/register`
                    break
                default:
                    throw new Error("Invalid backend")
            }

            // Create FormData to send files
            const formData = new FormData()
            formData.append("username", username)
            formData.append("password", password)
            formData.append("profilePhoto", profilePhoto)
            formData.append("idPhoto", idPhoto)

            const response = await fetch(apiUrl, {
                method: "POST",
                body: formData,
                // Don't set Content-Type header, it will be set automatically with boundary
            })

            const data = await response.json()
            const endTime = performance.now()
            const registerTimeMs = endTime - startTime
            const registerTime = `${Math.floor(registerTimeMs / 1000)}.${Math.floor(registerTimeMs % 1000)}` // seconds.milliseconds

            if (!response.ok) {
                throw new Error(data.message || "Registration failed")
            }

            toast("Registration successful", {
                description: `Registered using ${backend} backend in ${registerTime} seconds`,
            })

            // Reset form
            setUsername("")
            setPassword("")
            setProfilePhoto(null)
            setIdPhoto(null)
            setProfilePhotoPreview(null)
            setIdPhotoPreview(null)
            if (profilePhotoRef.current) profilePhotoRef.current.value = ""
            if (idPhotoRef.current) idPhotoRef.current.value = ""
        } catch (error) {
            toast("Registration failed", {
                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                // @ts-expect-error
                description: error.message,
            })
        } finally {
            setLoading(false)
        }
    }

    const renderTabContent = (backend: string) => (
        <>
            <CardContent className="space-y-4 pt-4">
                <div className="space-y-2">
                    <Label htmlFor={`username-${backend}`}>Username</Label>
                    <Input
                        id={`username-${backend}`}
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        placeholder="Enter your username"
                    />
                </div>
                <div className="space-y-2">
                    <Label htmlFor={`password-${backend}`}>Password</Label>
                    <Input
                        id={`password-${backend}`}
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        placeholder="Enter your password"
                    />
                </div>
                <div className="space-y-2">
                    <Label htmlFor={`profile-photo-${backend}`}>Profile Photo</Label>
                    <div className="flex items-center gap-4">
                        <Avatar className="h-16 w-16">
                            <AvatarImage src={profilePhotoPreview || ""} alt="Profile" />
                            <AvatarFallback>PH</AvatarFallback>
                        </Avatar>
                        <Input
                            id={`profile-photo-${backend}`}
                            type="file"
                            ref={profilePhotoRef}
                            accept="image/*"
                            onChange={handleProfilePhotoChange}
                        />
                    </div>
                </div>
                <div className="space-y-2">
                    <Label htmlFor={`id-photo-${backend}`}>ID Photo</Label>
                    <div className="flex items-center gap-4">
                        <div className="relative h-16 w-24 overflow-hidden rounded-md border">
                            {idPhotoPreview ? (
                                <img src={idPhotoPreview || "/placeholder.svg"} alt="ID" className="h-full w-full object-cover" />
                            ) : (
                                <div className="flex h-full w-full items-center justify-center bg-muted text-xs">No ID</div>
                            )}
                        </div>
                        <Input
                            id={`id-photo-${backend}`}
                            type="file"
                            ref={idPhotoRef}
                            accept="image/*"
                            onChange={handleIdPhotoChange}
                        />
                    </div>
                </div>
            </CardContent>
            <CardFooter className={"pt-5"}>
                <Button className="w-full" onClick={() => handleRegister(backend)} disabled={loading}>
                    {loading ? "Registering..." : `Register with ${backend}`}
                </Button>
            </CardFooter>
        </>
    )

    return (
        <Card>
            <Tabs defaultValue="axum">
                <TabsList className="grid w-full grid-cols-3">
                    <TabsTrigger value="axum">Axum</TabsTrigger>
                    <TabsTrigger value="rocket">Rocket</TabsTrigger>
                    <TabsTrigger value="gofiber">GoFiber</TabsTrigger>
                </TabsList>
                <TabsContent value="axum">{renderTabContent("axum")}</TabsContent>
                <TabsContent value="rocket">{renderTabContent("rocket")}</TabsContent>
                <TabsContent value="gofiber">{renderTabContent("gofiber")}</TabsContent>
            </Tabs>
        </Card>
    )
}
