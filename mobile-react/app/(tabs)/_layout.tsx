import React, { useEffect } from 'react';
import { Tabs } from "expo-router";
import TabBar from '@/components/TabBar';
import { useCartStore } from '@/store/cardBadgeStore';

export default function TabLayout() {
  const { fetchCartItems } = useCartStore();

  useEffect(() => {
    fetchCartItems();
  }, []);

  return (
    <Tabs tabBar={props => <TabBar {...props} />} screenOptions={{headerShown: false}}>
      <Tabs.Screen 
        name='index' 
        options={{
          tabBarLabel: 'Home', 
        }} 
      />
      <Tabs.Screen 
        name='explore' 
        options={{
          tabBarLabel: 'Explore',
        }} 
      />
      <Tabs.Screen 
        name='notifications' 
        options={{
          tabBarLabel: 'Alerts',
        }} 
      />
      <Tabs.Screen 
        name='cart' 
        options={{
          tabBarLabel: 'Cart',
        }} 
      />
      <Tabs.Screen 
        name='profile' 
        options={{
          tabBarLabel: 'Profile',
        }} 
      />
    </Tabs>
  );
}