import { FlatList, StyleSheet, Text, TouchableOpacity, View } from 'react-native'
import React from 'react'
import { useHeaderHeight } from '@react-navigation/elements'
import { Link, Stack } from 'expo-router'
import CartItem from '@/components/CartItem'
import { Colors } from '@/constants/Colors'
import Animated, { FadeInDown, SlideInDown } from 'react-native-reanimated'
import { useCartStore } from '@/store/cardBadgeStore'

type Props = {}

const CartScreen = () => {
  const headerHeight = useHeaderHeight()
  const { cartItems, fetchCartItems, total } = useCartStore()

  React.useEffect(() => {
    fetchCartItems()
  }, [])

  return (
    <>
      <Stack.Screen
        options={{ 
          headerShown: true, 
          headerTransparent: true, 
          headerTitleAlign: 'center', 
          title: 'Cart',
          headerRight: () => (
            <Link href="/cart" asChild>
              <TouchableOpacity style={{ marginRight: 20 }}>
                <Text style={{ color: Colors.primary, fontWeight: '500' }}>Total: R${total.toFixed(2)}</Text>
              </TouchableOpacity>
            </Link>
          )
        }}
      />
      <View style={[styles.container, { marginTop: headerHeight }]}>
        <FlatList
          data={cartItems}
          keyExtractor={(item) => item.productId.toString()}
          renderItem={({ item, index }) => (
            <Animated.View entering={FadeInDown.delay(300 + index * 100).duration(500)}>
              <CartItem item={item} />
            </Animated.View>
          )}
          ListEmptyComponent={
            <View style={styles.emptyContainer}>
              <Text style={styles.emptyText}>Your cart is empty</Text>
            </View>
          }
        />
      </View>
      {cartItems.length > 0 && (
        <Animated.View style={styles.footer} entering={SlideInDown.delay(500).duration(500)}>
          <Link href="/checkout" asChild>
            <TouchableOpacity style={styles.checkoutBtn}>
              <Text style={styles.checkoutBtnTxt}>Checkout R${total.toFixed(2)}</Text>
            </TouchableOpacity>
          </Link>
        </Animated.View>
      )}
    </>
  )
}

export default CartScreen

const styles = StyleSheet.create({
  container: {
    flex: 1,
    paddingHorizontal: 20
  },
  footer: {
    flexDirection: 'row',
    padding: 20,
    backgroundColor: Colors.white,
  },
  priceInfoWrapper: {
    flex: 1,
    justifyContent: 'center'
  },
  totalTxt: {
    fontSize: 16,
    fontWeight: '500',
    color: Colors.black
  },
  checkoutBtn: {
    flex: 1,
    backgroundColor: Colors.primary,
    height: 40,
    justifyContent: 'center',
    alignItems: 'center',
    borderRadius: 5
  },
  checkoutBtnTxt: {
    fontSize: 16,
    fontWeight: '500',
    color: Colors.white
  },
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    marginTop: 100
  },
  emptyText: {
    fontSize: 18,
    color: Colors.gray,
    fontWeight: '500'
  }
})