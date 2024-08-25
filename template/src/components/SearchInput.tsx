import { useState } from "react";

type Props = {
  setSearchAddress: (address: string) => void;
};

export const SearchInput: React.FC<Props> = ({ setSearchAddress }) => {
  const [word, setWord] = useState<string>("");

  const onChangeWord = (e: React.ChangeEvent<HTMLInputElement>) => {
    setWord(e.target.value);
  };

  const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setSearchAddress(word);
  };

  return (
    <>
      <form onSubmit={onSubmit}>
        <input
          type="text"
          className="py-2 px-3 ps-9 block w-full border-gray-200 shadow-sm rounded-lg text-sm focus:z-10 focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:pointer-events-none"
          value={word}
          onChange={onChangeWord}
          placeholder="Search for emails..."
        />
        <div className="absolute inset-y-0 start-0 flex items-center pointer-events-none ps-3">
          <svg
            className="size-4 text-gray-400"
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <circle cx="11" cy="11" r="8"></circle>
            <path d="m21 21-4.3-4.3"></path>
          </svg>
        </div>
      </form>
    </>
  );
};
