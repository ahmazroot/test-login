import Link from "next/link"
import { LoginForm } from "@/components/login-form"
import { Button } from "@/components/ui/button"

export default function Home() {
  return (
      <main className="flex min-h-screen flex-col items-center justify-center p-4 bg-gray-50">
        <div className="w-full max-w-md">
          <h1 className="text-2xl font-bold text-center mb-6">Login</h1>
          <LoginForm />

          <div className="mt-6 text-center">
            <p className="text-sm text-gray-600 mb-2">Don&#39;t have an account?</p>
            <Button variant="outline" asChild>
              <Link href="/registrasi">Register Now</Link>
            </Button>
          </div>
        </div>
      </main>
  )
}