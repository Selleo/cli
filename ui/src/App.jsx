import logo from './logo.svg';

function App() {
  return (
      <>
    <nav class="bg-pink-800">
        <div class="max-w-7xl mx-auto px-2 sm:px-6 lg:px-8">
            <div class="relative flex items-center justify-between h-16">
                <div class="absolute inset-y-0 left-0 flex items-center sm:hidden">
                    <button type="button"
                        class="inline-flex items-center justify-center p-2 rounded-md text-white hover:text-white hover:bg-pink-700 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white"
                        aria-controls="mobile-menu" aria-expanded="false">
                        <span class="sr-only">Open main menu</span>
                        <svg class="block h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                            stroke-width="2" stroke="currentColor" aria-hidden="true">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" />
                        </svg>
                        <svg class="hidden h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                            stroke-width="2" stroke="currentColor" aria-hidden="true">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>
                <div class="flex-1 flex items-center justify-center sm:items-stretch sm:justify-start">
                    <div class="flex-shrink-0 flex items-center">
                        <img class="block lg:hidden h-8 w-auto"
                            src="https://d33wubrfki0l68.cloudfront.net/9d77da3a3dd8cb65aa3f65cc257c37e4386c2440/3e547/selleo-logo-alt.svg" alt="Workflow"/>
                        <img class="hidden lg:block h-8 w-auto"
                            src="https://d33wubrfki0l68.cloudfront.net/9d77da3a3dd8cb65aa3f65cc257c37e4386c2440/3e547/selleo-logo-alt.svg"
                            alt="Selleo"/>
                    </div>
                    <div class="hidden sm:block sm:ml-6">
                        <div class="flex space-x-4">
                            <div class="relative inline-block text-left">
                                <div>
                                    <button type="button"
                                        class="inline-flex justify-center w-full rounded-md border border-pink-800 shadow-sm px-4 py-2 bg-pink-800 text-sm font-medium text-pink-100 hover:text-white hover:bg-pink-700"
                                        id="menu-button" aria-expanded="true" aria-haspopup="true">
                                        AWS PROFILE
                                        <svg class="-mr-1 ml-2 h-5 w-5" xmlns="http://www.w3.org/2000/svg"
                                            viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                                            <path fill-rule="evenodd"
                                                d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
                                                clip-rule="evenodd" />
                                        </svg>
                                    </button>
                                </div>

                                <div class="hidden origin-top-left absolute left-0 mt-2 w-56 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 focus:outline-none"
                                    role="menu" aria-orientation="vertical" aria-labelledby="menu-button" tabindex="-1">
                                    <div class="py-1" role="none">
                                        <a href="#" class="text-gray-700 block px-4 py-2 text-sm" role="menuitem"
                                            tabindex="-1" id="menu-item-0">Account settings</a>
                                        <a href="#" class="text-gray-700 block px-4 py-2 text-sm" role="menuitem"
                                            tabindex="-1" id="menu-item-1">Support</a>
                                        <a href="#" class="text-gray-700 block px-4 py-2 text-sm" role="menuitem"
                                            tabindex="-1" id="menu-item-2">License</a>
                                    </div>
                                </div>
                            </div>

                            <a href="#" class="bg-pink-900 text-white px-3 py-2 rounded-md text-sm font-medium"
                                aria-current="page">Secrets</a>

                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="sm:hidden" id="mobile-menu">
            <div class="px-2 pt-2 pb-3 space-y-1">
                <a href="#" class="bg-pink-900 text-white block px-3 py-2 rounded-md text-base font-medium"
                    aria-current="page">Secrets</a>

            </div>
        </div>
    </nav>





<div class="max-w-7xl mx-auto sm:px-6 lg:px-8 sm:py-6 lg:py-8">
<div class="bg-white overflow-hidden shadow rounded-lg divide-y divide-gray-200">
  <div class="px-4 py-5 sm:px-6">
<nav class="flex" aria-label="Breadcrumb">
  <ol role="list" class="flex items-center space-x-2">
    <li>
      <div class="flex items-center">
        <svg class="flex-shrink-0 h-5 w-5 text-gray-300" xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
          <path d="M5.555 17.776l8-16 .894.448-8 16-.894-.448z" />
        </svg>
        <a href="#" class="ml-2 text-sm font-medium text-pink-700 hover:text-pink-500">portal</a>
      </div>
    </li>
    <li>
      <div class="flex items-center">
        <svg class="flex-shrink-0 h-5 w-5 text-gray-300" xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
          <path d="M5.555 17.776l8-16 .894.448-8 16-.894-.448z" />
        </svg>
        <a href="#" class="ml-2 text-sm font-medium text-pink-700 hover:text-pink-500">development</a>
      </div>
    </li>

    <li>
      <div class="flex items-center">
        <svg class="flex-shrink-0 h-5 w-5 text-gray-300" xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
          <path d="M5.555 17.776l8-16 .894.448-8 16-.894-.448z" />
        </svg>
        <a href="#" class="ml-2 text-sm font-medium text-gray-500 hover:text-pink-500" aria-current="page">api</a>
      </div>
    </li>
  </ol>
</nav>
  </div>
  <div class="px-4 py-5 sm:p-6">

<form class="space-y-8 divide-y divide-gray-200" autocomplete="off">
  <div class="space-y-8 divide-y divide-gray-200 sm:space-y-5">
    <div>
      <div>
        <h3 class="text-lg leading-6 font-medium text-gray-900">Secrets</h3>
        <p class="mt-1 max-w-2xl text-sm text-gray-500">Last fetched at 14 June 2022 at 00:54.</p>
      </div>

      <div class="mt-6 sm:mt-5 space-y-6 sm:space-y-5">
        <div class="sm:grid sm:grid-cols-4 sm:gap-4 sm:items-start sm:border-t sm:border-gray-200 sm:pt-5">
          <label for="env-0" class="block text-sm font-medium text-gray-700 sm:mt-px sm:pt-2">DATABASE_URL</label>
          <div class="mt-1 sm:mt-0 sm:col-span-3">
            <input id="env-0" name="env-0" type="text" autocomplete="off" readonly disabled class="appearance-none block w-full shadow-sm bg-gray-100 focus:ring-pink-500 focus:border-pink-500 sm:text-sm border-gray-300 rounded-md"
             value="postgres://aaaaaaaaaaaaa:bbbbbbbbbbbbbb@ec2-10-0-0-200.eu-west-1.compute.amazonaws.com:5432/api"/>
          </div>
        </div>

        <div class="sm:grid sm:grid-cols-4 sm:gap-4 sm:items-start sm:border-t sm:border-gray-200 sm:pt-5">
          <label for="env-1" class="block text-sm font-medium text-gray-700 sm:mt-px sm:pt-2">SECRET_KEY_BASE</label>
          <div class="mt-1 sm:mt-0 sm:col-span-3">
            <input id="env-1" name="env-1" type="text" autocomplete="off" readonly disabled class="appearance-none block w-full shadow-sm bg-gray-100 focus:ring-pink-500 focus:border-pink-500 sm:text-sm border-gray-300 rounded-md"
            value="41cdbba4015e7ca687674f995f2fa37402fd456ca38ad4bb84344102b2439d686f2b49750f9cca015378fa793fdfcf98b8d86145ce35bc5ad77c8be11cbb1807"/>
          </div>
        </div>

        <div class="sm:grid sm:grid-cols-4 sm:gap-4 sm:items-start sm:border-t sm:border-gray-200 sm:pt-5">
          <label for="env-2" class="block text-sm font-medium text-gray-700 sm:mt-px sm:pt-2">SLACK_WEBHOOK_URL</label>
          <div class="mt-1 sm:mt-0 sm:col-span-3">
            <input id="env-2" name="env-2" type="text" autocomplete="off" class="appearance-none block w-full shadow-sm focus:ring-pink-500 focus:border-pink-500 sm:text-sm border-gray-300 rounded-md"
value="https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
            />
          </div>
        </div>

      </div>
    </div>

  </div>

  <div class="pt-5">
    <div class="flex justify-end">
      <button type="submit" class="ml-3 inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-pink-600 hover:bg-pink-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-pink-500">Save changes</button>
    </div>
  </div>
</form>
  </div>
</div>
</div>
  </>
  );
}

export default App;
