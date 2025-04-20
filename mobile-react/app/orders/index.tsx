import { View, Text, FlatList, StyleSheet, TouchableOpacity, ActivityIndicator } from 'react-native'
import React, { useEffect } from 'react'
import { Stack, useRouter } from 'expo-router'
import { useHeaderHeight } from '@react-navigation/elements'
import { useOrderStore } from '@/store/orderStore'
import { Colors } from '@/constants/Colors'
import { OrderStatus } from '@/service/orderService'
import Animated, { FadeInDown } from 'react-native-reanimated'
import Ionicons from 'react-native-vector-icons/Ionicons'

const StatusColors = {
  warning: '#FFA000',
  info: '#2196F3',
  success: '#4CAF50',
  error: '#F44336'
}

const OrdersScreen = () => {
  const router = useRouter()
  const headerHeight = useHeaderHeight()
  const { fetchOrders, orders, isLoading, error, clearError } = useOrderStore()

  useEffect(() => {
    fetchOrders()
  }, [])

  const getStatusColor = (status: OrderStatus) => {
    switch (status) {
      case OrderStatus.PENDING:
        return StatusColors.warning
      case OrderStatus.PROCESSING:
        return StatusColors.info
      case OrderStatus.SHIPPED:
        return '#1E88E5' 
      case OrderStatus.DELIVERED:
        return StatusColors.success
      case OrderStatus.CANCELLED:
        return StatusColors.error
      default:
        return Colors.gray
    }
  }

  const getStatusLabel = (status: OrderStatus) => {
    switch (status) {
      case OrderStatus.PENDING:
        return 'Pending'
      case OrderStatus.PROCESSING:
        return 'Processing'
      case OrderStatus.SHIPPED:
        return 'Shipped'
      case OrderStatus.DELIVERED:
        return 'Delivered'
      case OrderStatus.CANCELLED:
        return 'Cancelled'
      default:
        return status
    }
  }

  if (isLoading) {
    return (
      <>
        <Stack.Screen
          options={{
            headerShown: true,
            headerTransparent: true,
            headerTitleAlign: 'center',
            title: 'My Orders',
          }}
        />
        <View style={[styles.container, { marginTop: headerHeight }, styles.centerContent]}>
          <ActivityIndicator size="large" color={Colors.primary} />
          <Text style={styles.loadingText}>Loading your orders...</Text>
        </View>
      </>
    )
  }

  return (
    <>
      <Stack.Screen
        options={{
          headerShown: true,
          headerTransparent: true,
          headerTitleAlign: 'center',
          title: 'My Orders',
        }}
      />

      <View style={[styles.container, { marginTop: headerHeight }]}>
        {error ? (
          <View style={styles.errorContainer}>
            <Ionicons name="alert-circle-outline" size={50} color={StatusColors.error} />
            <Text style={styles.errorText}>{error}</Text>
            <TouchableOpacity
              style={styles.retryButton}
              onPress={() => {
                clearError()
                fetchOrders()
              }}
            >
              <Text style={styles.retryButtonText}>Try Again</Text>
            </TouchableOpacity>
          </View>
        ) : orders.length === 0 ? (
          <View style={styles.emptyContainer}>
            <Ionicons name="cart-outline" size={60} color={Colors.gray} />
            <Text style={styles.emptyText}>You don't have any orders yet</Text>
            <TouchableOpacity
              style={styles.shopButton}
              onPress={() => router.push("/(tabs)")}
            >
              <Text style={styles.shopButtonText}>Go Shopping</Text>
            </TouchableOpacity>
          </View>
        ) : (
          <FlatList
            data={orders}
            keyExtractor={(item) => item.id.toString()}
            renderItem={({ item, index }) => (
              <Animated.View entering={FadeInDown.delay(100 * index).duration(400)}>
                <TouchableOpacity
                  style={styles.orderCard}
                  onPress={() => router.push({
                    pathname: "/orders/[id]",
                    params: { id: item.id }
                  })}
                >
                  <View style={styles.orderHeader}>
                    <Text style={styles.orderNumber}>Order #{item.id}</Text>
                    <View style={[styles.statusBadge, { backgroundColor: getStatusColor(item.status) }]}>
                      <Text style={styles.statusText}>{getStatusLabel(item.status)}</Text>
                    </View>
                  </View>
                  
                  <View style={styles.orderInfo}>
                    <Text style={styles.orderDate}>{item.createdAt}</Text>
                    <Text style={styles.orderTotal}>R${item.totalAmount.toFixed(2)}</Text>
                  </View>
                  
                  <View style={styles.orderFooter}>
                    <Text style={styles.viewDetails}>View Details</Text>
                    <Ionicons name="chevron-forward" size={20} color={Colors.primary} />
                  </View>
                </TouchableOpacity>
              </Animated.View>
            )}
            contentContainerStyle={styles.listContainer}
          />
        )}
      </View>
    </>
  )
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: Colors.background,
    padding: 16,
  },
  centerContent: {
    justifyContent: 'center',
    alignItems: 'center',
  },
  loadingText: {
    marginTop: 16,
    fontSize: 16,
    color: Colors.gray,
  },
  errorContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  errorText: {
    marginTop: 10,
    marginBottom: 20,
    fontSize: 16,
    color: StatusColors.error,
    textAlign: 'center',
  },
  retryButton: {
    backgroundColor: Colors.primary,
    paddingVertical: 10,
    paddingHorizontal: 20,
    borderRadius: 8,
  },
  retryButtonText: {
    color: Colors.white,
    fontSize: 16,
    fontWeight: '600',
  },
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  emptyText: {
    marginTop: 16,
    marginBottom: 24,
    fontSize: 18,
    color: Colors.gray,
    textAlign: 'center',
  },
  shopButton: {
    backgroundColor: Colors.primary,
    paddingVertical: 12,
    paddingHorizontal: 24,
    borderRadius: 8,
  },
  shopButtonText: {
    color: Colors.white,
    fontSize: 16,
    fontWeight: '600',
  },
  listContainer: {
    paddingBottom: 20,
  },
  orderCard: {
    backgroundColor: Colors.white,
    borderRadius: 12,
    padding: 16,
    marginVertical: 8,
    shadowColor: Colors.darkGray,
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 2,
  },
  orderHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  orderNumber: {
    fontSize: 16,
    fontWeight: '700',
    color: Colors.black,
  },
  statusBadge: {
    paddingVertical: 4,
    paddingHorizontal: 10,
    borderRadius: 20,
  },
  statusText: {
    color: Colors.white,
    fontSize: 12,
    fontWeight: '600',
  },
  orderInfo: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 12,
    paddingBottom: 12,
    borderBottomWidth: 1,
    borderBottomColor: Colors.extraLightGray,
  },
  orderDate: {
    fontSize: 14,
    color: Colors.darkGray,
  },
  orderTotal: {
    fontSize: 15,
    fontWeight: '600',
    color: Colors.black,
  },
  orderFooter: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'flex-end',
  },
  viewDetails: {
    fontSize: 14,
    color: Colors.primary,
    marginRight: 4,
  },
})

export default OrdersScreen 