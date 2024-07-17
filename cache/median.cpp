#include <iostream>
#include <vector>
using namespace std;

double findMedianSortedArrays(vector<int>& nums1, vector<int>& nums2) {
  double median;
  int i=0,j=0,m,n,mid,k=0;
  m= nums1.size();
  n= nums2.size();
  if((m+n)%2 != 0){
    mid = (m+n+1)/2;
    while(k<mid){
      if(nums1[i]<nums2[j] && i< (m-1)){
        i++;
      }else if(nums1[i]>nums2[j] && j<(n-1)){
        j++;
      }else{
        break;
      }
      k++;
    }
    if(nums1[i]>nums2[j] && i<k){
      return nums1[i];
    }else{
      return nums2[j];
    }
  }else{
    mid = (m+n)/2;
    while(k<mid){
      if(nums1[i]<nums2[j] && i< (m-1)){
        i++;
      }else if(nums1[i]>nums2[j] && j<(n-1)){
        j++;
      }else{
        break;
      }
      k++;
    }
    if(nums1[i]<nums2[j]){
      if(nums1[i+1]<nums2[j] && i+1 <m){
        return ((nums1[i]+nums1[i+1])/(double)2);
      }else{
        return ((nums1[i]+nums2[j])/(double)2);
      }
    }else{
      if(nums2[j+1]<nums1[i] && j+1<n){
        return ((nums2[j]+nums2[j+1])/(double)2);
      }else{
        return ((nums1[i]+nums2[j])/(double)2);
      }
    }
  }
}

int main(){
  vector<int> num1,num2;
  num1 = {1,3};
  num2 = {2};
  cout << findMedianSortedArrays(num1,num2)<<endl;
}