"use client";

import { useState } from "react";
import { faEnvelope, faEye, faEyeSlash, faLock } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

export default function Home() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [passwordVisible, setPasswordVisible] = useState(false);

  const handleSubmit = (e: any) => {
    e.preventDefault();
    console.log("Email:", email);
    console.log("Password:", password);
  };

  const togglePasswordVisibility = () => {
    setPasswordVisible(!passwordVisible);
  };

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24">
      {/* Quadrado centralizado */}
      <div className="bg-white p-8 rounded-lg shadow-lg">
        {/* Formulário com campos de email e senha */}
        <form onSubmit={handleSubmit} className="flex flex-col space-y-4">
          <div>
            {/* Campo de email com ícone e placeholder */}
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
            {/* Campo de senha com ícone e placeholder */}
            <div className="relative">
              <input
                type="password"
                id="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Senha"
                className="border rounded p-2 w-full pl-8" // Adicione espaço de preenchimento à esquerda
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
          {/* Botão Entrar com largura total */}
          <button type="submit" className="w-full bg-blue-500 text-white p-2 rounded">
            Entrar
          </button>
          {/* Links centralizados abaixo do botão Entrar */}
          <div className="flex justify-center mt-4 space-x-4 w-full">
            <a href="/recover-password" className="text-sm text-blue-500 w-1/2 text-center">
              Esqueci minha senha
            </a>
            <a href="/sign-up" className="text-sm text-blue-500 w-1/2 text-center">
              Não possuo uma conta
            </a>
          </div>
        </form>
      </div>
    </main>
  );
}
