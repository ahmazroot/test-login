import { RegisterForm } from "@/components/register-form"

export default function RegisterPage() {
    return (
        <main className="flex min-h-screen flex-col items-center justify-center p-4 bg-gray-50">
            <div className="w-full max-w-md">
                <h1 className="text-2xl font-bold text-center mb-6">register</h1>
                <RegisterForm />
            </div>
        </main>
    )
}
