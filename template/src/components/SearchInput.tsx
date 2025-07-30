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
    <form onSubmit={onSubmit} className="w-full max-w-lg">
      <label
        htmlFor="email-search"
        className="mb-2 text-sm font-medium text-gray-900 sr-only"
      >
        メールアドレスで検索
      </label>
      <div className="relative">
        <div className="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
          <svg
            className="w-5 h-5 text-gray-500"
            aria-hidden="true"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 20 20"
          >
            <path
              stroke="currentColor"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"
            />
          </svg>
        </div>
        <input
          type="search"
          id="email-search"
          className="block w-full p-4 ps-10 text-sm text-gray-900 border border-gray-300 rounded-lg bg-white shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-shadow duration-200"
          value={word}
          onChange={onChangeWord}
          placeholder="メールアドレスを入力して検索..."
          required
        />
      </div>
    </form>
  );
};
