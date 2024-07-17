#include <iostream>
using namespace std;

int gcd(int a, int b)
{
    int R;
    while ((a % b) > 0)
    {
        R = a % b;
        a = b;
        b = R;
    }
    return b;
}

int main()
{
    int arr[9] = {1, 2, 3, 4, 5, 6, 7, 8, 9}, n = 9, k = 3;
    int number_of_cycle = gcd(n, k), cycle_length = n / number_of_cycle;
    int temp;
    for (int i = 0; i < number_of_cycle; i++)
    {
        temp = arr[i];
        int j = i;
        cout << "Cycle " << i + 1 << " :\n";
        for (int c = 0; c < cycle_length - 1; c++)
        {
            cout << j << "-->";
            arr[j] = arr[(j + k) % n];
            j = (j + k) % n;
        }
        arr[j] = temp;
        cout <<j<< endl;
    }
    cout << "Sequence of elements:\n";
    for (int i = 0; i < n; i++)
    {
        cout << arr[i] << endl;
    }

    return 0;
}