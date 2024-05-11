import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faUser } from "@fortawesome/free-solid-svg-icons";
import { FieldError, UseFormRegisterReturn } from "react-hook-form";

type InputNameProps = {
    name: UseFormRegisterReturn; // Atualize a prop para aceitar UseFormRegisterReturn
    errors?: FieldError | undefined;
};

export function InputName({ name, errors }: InputNameProps) {
    return (
        <div className="relative">
            <input
                id="name"
                type="text"
                {...name} // Use o objeto de registro retornado por register
                placeholder="Name"
                className="border rounded p-2 w-full pl-8"
            />
            <div className="absolute inset-y-0 left-0 flex items-center pl-2">
                <FontAwesomeIcon icon={faUser} className="text-gray-400" />
            </div>
            {errors && (
                <span className="text-red-500 text-sm mt-1">{errors.message}</span>
            )}
        </div>
    );
}
