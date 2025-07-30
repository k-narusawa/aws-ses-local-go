export const TableHeader: React.FC = () => {
  return (
    <thead className="bg-gray-50">
      <tr>
        <th scope="col" className="py-3 px-4 pe-0">
          <div className="flex items-center h-5">
            <input
              id="select-all-emails"
              type="checkbox"
              className="w-4 h-4 border-gray-300 rounded text-blue-600 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors duration-200"
              aria-label="全てのメールを選択"
            />
          </div>
        </th>
        <th
          scope="col"
          className="px-6 py-3 text-start text-xs font-medium text-gray-600 uppercase tracking-wider"
        >
          件名
        </th>
        <th
          scope="col"
          className="px-6 py-3 text-start text-xs font-medium text-gray-600 uppercase tracking-wider"
        >
          送信元
        </th>
        <th
          scope="col"
          className="px-6 py-3 text-start text-xs font-medium text-gray-600 uppercase tracking-wider"
        >
          宛先
        </th>
        <th
          scope="col"
          className="px-6 py-3 text-end text-xs font-medium text-gray-600 uppercase tracking-wider"
        >
          操作
        </th>
      </tr>
    </thead>
  );
};
