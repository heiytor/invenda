import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faEye, faEyeSlash, faLock } from "@fortawesome/free-solid-svg-icons";
import { useState } from "react";
import { FieldError, UseFormRegisterReturn } from "react-hook-form";

type InputPasswordProps = {
    id?: string;
    password: UseFormRegisterReturn;
    errors?: FieldError | undefined;
    placeholder?: string;
};

export function InputPassword({ password, errors, id="password", placeholder="Senha" }: InputPasswordProps) {
    const [passwordVisible, setPasswordVisible] = useState(false);
    const togglePassword = () => {
        setPasswordVisible(!passwordVisible);
    };

    return (
        <div className="relative">
            <input
                id={id}
                type={passwordVisible ? "text" : "password"}
                {...password}
                placeholder={placeholder}
                className="border rounded p-2 w-full pl-8"
            />
            <div className="absolute inset-y-0 left-0 flex items-center pl-2">
                <FontAwesomeIcon icon={faLock} className="text-gray-400" />
            </div>
            <div className="absolute inset-y-0 right-0 flex items-center pr-2 cursor-pointer" onClick={togglePassword}>
                <FontAwesomeIcon icon={passwordVisible ? faEye : faEyeSlash} className="text-gray-400" />
            </div>
            {errors && (
                <span className="text-red-500 text-sm mt-1">{errors.message}</span>
            )}
        </div>
    );
}
