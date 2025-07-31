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
              aria-label="Select all emails"
            />
          </div>
        </th>
        <th
          scope="col"
          className="px-6 py-3 text-start text-xs font-medium text-gray-600 uppercase tracking-wider"
        >
          Subject
        </th>
        <th
          scope="col"
          className="px-6 py-3 text-start text-xs font-medium text-gray-600 uppercase tracking-wider"
        >
          From
        </th>
        <th
          scope="col"
          className="px-6 py-3 text-start text-xs font-medium text-gray-600 uppercase tracking-wider"
        >
          To
        </th>
        <th
          scope="col"
          className="px-6 py-3 text-end text-xs font-medium text-gray-600 uppercase tracking-wider"
        >
          Actions
        </th>
      </tr>
    </thead>
  );
};
