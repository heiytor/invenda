"use client";

import { InputEmail } from "@/components/Inputs/Email";
import { InputName } from "@/components/Inputs/Name";
import { InputPassword } from "@/components/Inputs/Password";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

const createUserSchema = z.object(
    {
        name: z.string().min(1, "TODO"),
        email: z.string().min(1, "TODO").email("TODO"),
        password: z.string().regex(/^(?=.*[A-Z])(?=.*[a-z])(?=.*\d)(?=.*[!@#$%^&*()_+\-=[\]{}|;:'",.<>?]).{8,64}$/, "TODO"),
        confirmPassword: z.string().min(1, "TODO"),
    },
).refine(({ password, confirmPassword}) => password === confirmPassword, {
  message: "Password doesn't match",
  path: ["confirmPassword"]
});

type CreateUserForm = z.infer<typeof createUserSchema>;

async function onSubmit(_: CreateUserForm): Promise<void> {
    // const { email, password } = data;
    // await authUser({ identifier: email, password: password });
}

export default function SignUp() {
    const { 
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<CreateUserForm>({ resolver: zodResolver(createUserSchema) });

    return (
        <main className="flex min-h-screen flex-col items-center justify-center p-24">
            <div className="bg-white p-8 rounded-lg shadow-lg">
                <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col space-y-4">
                    <InputName name={register("email")} errors={errors.email}/>
                    <InputEmail email={register("email")} errors={errors.email}/>
                    <InputPassword id="password" password={register("password")} errors={errors.password}/>
                    <InputPassword id="confirmPassword" password={register("confirmPassword")} errors={errors.confirmPassword} placeholder="Confirmar senha"/>

                    <button type="submit" className="w-full bg-blue-500 text-white p-2 rounded">
                        Cadastrar
                    </button>
                    <div className="flex justify-center mt-4 space-x-4 w-full">
                        <a href="/sign-in" className="text-sm text-blue-500 w-full text-center">
                            Já possuo uma conta.
                        </a>
                    </div>
                </form>
            </div>
        </main>
    );
}
