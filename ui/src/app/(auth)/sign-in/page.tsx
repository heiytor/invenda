"use client";

import { InputEmail } from "@/components/Inputs/Email";
import { InputPassword } from "@/components/Inputs/Password";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { authUser } from "@/api/auth";

const authUserSchema = z.object({
    email: z.string().min(1, "TODO").email("Email inválido"),
    password: z.string().min(1, "TODO"),
    // password: z.string().regex(/^(?=.*[A-Z])(?=.*[a-z])(?=.*\d)(?=.*[!@#$%^&*()_+\-=[\]{}|;:'",.<>?]).{8,64}$/, "A senha deve ter entre 8 e 64 caracteres, incluindo pelo menos uma letra maiúscula, uma letra minúscula, um número e um caractere especial."),
});

type AuthUserForm = z.infer<typeof authUserSchema>;

async function onSubmit(data: AuthUserForm): Promise<void> {
    const { email, password } = data;
    await authUser({ identifier: email, password: password });
}

export default function Home() {
    const { 
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<AuthUserForm>({ resolver: zodResolver(authUserSchema) });

    return (
        <main className="flex min-h-screen flex-col items-center justify-center p-24">
            <div className="bg-white p-8 rounded-lg shadow-lg">
                <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col space-y-4">
                    <InputEmail email={register("email")} errors={errors.email}/>
                    <InputPassword password={register("password")} errors={errors.password}/>
                    <button type="submit" className="w-full bg-blue-500 text-white p-2 rounded"> Entrar </button>
                    <div className="flex justify-center mt-4 space-x-4 w-full">
                        <a href="/recover-password" className="text-sm text-blue-500 w-1/2 text-center">
                            Esqueci minha senha
                        </a>
                        <a href="/sign-up" className="text-sm text-blue-500 w-1/2 text-center">
                            Cadastrar
                        </a>
                    </div>
                </form>
            </div>
        </main>
    );
}
