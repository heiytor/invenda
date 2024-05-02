"use client";

import { useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faEye, faEyeSlash, faUser, faLock, faEnvelope } from "@fortawesome/free-solid-svg-icons";

export default function SignUp() {
    const [name, setName] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [passwordVisible, setPasswordVisible] = useState(false);
    const [confirmPasswordVisible, setConfirmPasswordVisible] = useState(false);
    const [termsAccepted, setTermsAccepted] = useState(false);
    const [receiveNotifications, setReceiveNotifications] = useState(false);

    const handleSubmit = (e: any) => {
        e.preventDefault();

        if (!termsAccepted) {
            console.log("Você precisa aceitar os termos de uso.");
            return;
        }

        if (password !== confirmPassword) {
            console.log("As senhas não coincidem.");
            return;
        }

        console.log('Name:', name);
        console.log('Email:', email);
        console.log('Password:', password);
        console.log('Confirm Password:', confirmPassword);
        console.log('Receive Notifications:', receiveNotifications);
        console.log('Terms Accepted:', termsAccepted);
    };

    const togglePasswordVisibility = () => {
        setPasswordVisible(!passwordVisible);
    };

    const toggleConfirmPasswordVisibility = () => {
        setConfirmPasswordVisible(!confirmPasswordVisible);
    };

    return (
        <main className="flex min-h-screen flex-col items-center justify-center p-24">
            {/* Quadrado centralizado */}
            <div className="bg-white p-8 rounded-lg shadow-lg">
                <form onSubmit={handleSubmit} className="flex flex-col space-y-4">
                    <div>
                        <div className="relative">
                            <input
                                type="text"
                                id="name"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                                placeholder="Nome"
                                className="border rounded p-2 w-full pl-8"
                                required
                            />
                            {/* Ícone de usuário */}
                            <div className="absolute inset-y-0 left-0 flex items-center pl-2">
                                <FontAwesomeIcon icon={faUser} className="text-gray-400" />
                            </div>
                        </div>
                    </div>
                    <div>
                        <div className="relative">
                            <input
                                type="email"
                                id="email"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                placeholder="Email"
                                className="border rounded p-2 w-full pl-8"
                                required
                            />
                            {/* Ícone de usuário */}
                            <div className="absolute inset-y-0 left-0 flex items-center pl-2">
                                <FontAwesomeIcon icon={faEnvelope} className="text-gray-400" />
                            </div>
                        </div>
                    </div>
                    <div>
                        <div className="relative">
                            <input
                                type={passwordVisible ? "text" : "password"}
                                id="password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                placeholder="Senha"
                                className="border rounded p-2 w-full pl-8"
                                required
                            />
                            {/* Ícone de cadeado */}
                            <div className="absolute inset-y-0 left-0 flex items-center pl-2">
                                <FontAwesomeIcon icon={faLock} className="text-gray-400" />
                            </div>
                            {/* Botão para visualizar a senha */}
                            <div className="absolute inset-y-0 right-0 flex items-center pr-2 cursor-pointer" onClick={togglePasswordVisibility}>
                                <FontAwesomeIcon icon={passwordVisible ? faEyeSlash : faEye} className="text-gray-400" />
                            </div>
                        </div>
                    </div>
                    <div>
                        <div className="relative">
                            <input
                                type={confirmPasswordVisible ? "text" : "password"}
                                id="confirmPassword"
                                value={confirmPassword}
                                onChange={(e) => setConfirmPassword(e.target.value)}
                                placeholder="Confirme a senha"
                                className="border rounded p-2 w-full pl-8"
                                required
                            />
                            {/* Ícone de cadeado */}
                            <div className="absolute inset-y-0 left-0 flex items-center pl-2">
                                <FontAwesomeIcon icon={faLock} className="text-gray-400" />
                            </div>
                            {/* Botão para visualizar a confirmação de senha */}
                            <div className="absolute inset-y-0 right-0 flex items-center pr-2 cursor-pointer" onClick={toggleConfirmPasswordVisibility}>
                                <FontAwesomeIcon icon={confirmPasswordVisible ? faEyeSlash : faEye} className="text-gray-400" />
                            </div>
                        </div>
                    </div>
                    {/* Checkbox para aceitar os termos de uso */}
                    <div className="flex items-center">
                        <input
                            type="checkbox"
                            id="termsAccepted"
                            checked={termsAccepted}
                            onChange={(e) => setTermsAccepted(e.target.checked)}
                            className="mr-2"
                            required
                        />
                        <label htmlFor="termsAccepted" className="text-sm">
                            Eu aceito os
                            {' '}
                            <a href="/termos-de-uso" className="text-blue-500 underline" target="_blank" rel="noopener noreferrer">
                            termos de uso
                            </a>.
                        </label>
                    </div>
                    {/* Checkbox para receber notificações por email */}
                    <div className="flex items-center">
                        <input
                            type="checkbox"
                            id="receiveNotifications"
                            checked={receiveNotifications}
                            onChange={(e) => setReceiveNotifications(e.target.checked)}
                            className="mr-2"
                        />
                        <label htmlFor="receiveNotifications" className="text-sm">
                            Receber notificações por e-mail.
                        </label>
                    </div>
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
