#include <iostream>
using namespace std;

void _swap(int *a,int *b){
    int temp = *a;
    *a = *b;
    *b = temp;
}

void _reverse(int arr[],int l,int r){
    int n = r-1+l;
    for(int i=l;i<n/2;i++){
        _swap(&arr[i],&arr[n-i-1]); 
    }
}

int main(){
    int arr[10]={1,2,3,4,5,6,7,8},n=8,j,swap;
    // for(int i=0;i<(n/2.0)-1;i++){
    //     j = n-i-1;
    //     swap = arr[i];
    //     arr[i] = arr[j];
    //     arr[j] = swap; 
    // }
    _reverse(arr,0,n);
    for(int i=0;i<n;i++){
        cout<<arr[i]<<endl;
    }

    return 0;
}