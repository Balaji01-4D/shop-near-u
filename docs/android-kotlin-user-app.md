# Android Kotlin User App Specification

This document provides a comprehensive guide for building an Android app in Kotlin for the user/customer side of the shop-near-u platform. Users can discover nearby shops, browse products, and manage their accounts.

## Backend API Overview

Based on the backend code analysis, the user-facing endpoints are:

### Authentication Endpoints
- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `GET /auth/me` - Get current user profile (requires auth)
- `POST /auth/logout` - Logout user (requires auth)
- `POST /auth/change-password` - Change password (requires auth)
- `DELETE /auth/delete-account` - Delete user account (requires auth)

### Shop Discovery Endpoints
- `GET /shop?lat={lat}&lon={lon}&radius={radius}&limit={limit}` - Find nearby shops

### Product Catalog Endpoints
- `GET /api/catalog-products/suggest?keyword={keyword}&limit={limit}` - Search products

## Architecture Overview

### Recommended Architecture: MVVM + Repository Pattern

```
┌─────────────────────────────────────────┐
│                   UI Layer               │
│  ┌─────────────┐  ┌─────────────────────┐│
│  │  Activities │  │     Fragments       ││
│  │  (Compose)  │  │   (Compose UI)      ││
│  └─────────────┘  └─────────────────────┘│
└─────────────────────────────────────────┘
                    │
┌─────────────────────────────────────────┐
│              ViewModel Layer             │
│  ┌─────────────┐  ┌─────────────────────┐│
│  │    Auth     │  │       Shop          ││
│  │  ViewModel  │  │    ViewModel        ││
│  └─────────────┘  └─────────────────────┘│
└─────────────────────────────────────────┘
                    │
┌─────────────────────────────────────────┐
│             Repository Layer             │
│  ┌─────────────┐  ┌─────────────────────┐│
│  │    Auth     │  │       Shop          ││
│  │ Repository  │  │   Repository        ││
│  └─────────────┘  └─────────────────────┘│
└─────────────────────────────────────────┘
                    │
┌─────────────────────────────────────────┐
│            Data Source Layer             │
│  ┌─────────────┐  ┌─────────────────────┐│
│  │   Remote    │  │       Local         ││
│  │ Data Source │  │   Data Source       ││
│  │  (Retrofit) │  │     (Room)          ││
│  └─────────────┘  └─────────────────────┘│
└─────────────────────────────────────────┘
```

## Key Dependencies

Add these to your `app/build.gradle.kts`:

```kotlin
dependencies {
    // Core Android
    implementation("androidx.core:core-ktx:1.12.0")
    implementation("androidx.lifecycle:lifecycle-runtime-ktx:2.7.0")
    implementation("androidx.activity:activity-compose:1.8.2")
    
    // Compose
    implementation(platform("androidx.compose:compose-bom:2023.10.01"))
    implementation("androidx.compose.ui:ui")
    implementation("androidx.compose.ui:ui-tooling-preview")
    implementation("androidx.compose.material3:material3")
    
    // Navigation
    implementation("androidx.navigation:navigation-compose:2.7.5")
    
    // ViewModel
    implementation("androidx.lifecycle:lifecycle-viewmodel-compose:2.7.0")
    
    // Networking
    implementation("com.squareup.retrofit2:retrofit:2.9.0")
    implementation("com.squareup.retrofit2:converter-gson:2.9.0")
    implementation("com.squareup.okhttp3:okhttp:4.12.0")
    implementation("com.squareup.okhttp3:logging-interceptor:4.12.0")
    
    // Image Loading
    implementation("io.coil-kt:coil-compose:2.5.0")
    
    // Local Storage
    implementation("androidx.room:room-runtime:2.6.1")
    implementation("androidx.room:room-ktx:2.6.1")
    kapt("androidx.room:room-compiler:2.6.1")
    
    // Location Services
    implementation("com.google.android.gms:play-services-location:21.0.1")
    implementation("com.google.android.gms:play-services-maps:18.2.0")
    
    // Dependency Injection
    implementation("io.insert-koin:koin-android:3.5.0")
    implementation("io.insert-koin:koin-androidx-compose:3.5.0")
    
    // Preferences
    implementation("androidx.datastore:datastore-preferences:1.0.0")
}
```

## Data Models

Create these Kotlin data classes in `data/models/`:

```kotlin
// data/models/ApiResponse.kt
data class ApiResponse<T>(
    val success: Boolean,
    val message: String,
    val data: T? = null
)

data class ApiError(
    val success: Boolean = false,
    val message: String,
    val error: String? = null
)

// data/models/User.kt
data class User(
    val id: Int,
    val name: String,
    val email: String,
    val latitude: Double = 0.0,
    val longitude: Double = 0.0,
    val role: String = "user",
    val createdAt: String? = null
)

data class UserRegisterRequest(
    val name: String,
    val email: String,
    val password: String
)

data class UserLoginRequest(
    val email: String,
    val password: String
)

data class ChangePasswordRequest(
    val oldPassword: String,
    val newPassword: String
)

data class UserAuthResponse(
    val user: User,
    val token: String
)

// data/models/Shop.kt
data class Shop(
    val id: Int,
    val name: String,
    val ownerName: String,
    val type: String,
    val email: String,
    val mobile: String,
    val address: String,
    val latitude: Double,
    val longitude: Double,
    val supportsDelivery: Boolean = false,
    val distance: Double? = null,
    val createdAt: String? = null
)

// data/models/Product.kt
data class CatalogProduct(
    val id: Int,
    val name: String,
    val brand: String? = null,
    val category: String,
    val description: String? = null,
    val imageUrl: String? = null,
    val createdAt: String? = null
)

data class ShopProduct(
    val id: Int,
    val shopId: Int,
    val catalogId: Int,
    val price: Double,
    val stock: Int,
    val isAvailable: Boolean,
    val discount: Double = 0.0,
    val createdAt: String? = null,
    val updatedAt: String? = null,
    val catalogProduct: CatalogProduct? = null
)
```

## Network Layer

### API Service Interface

```kotlin
// network/ApiService.kt
import retrofit2.Response
import retrofit2.http.*

interface ApiService {
    
    // Authentication
    @POST("auth/register")
    suspend fun register(@Body request: UserRegisterRequest): Response<ApiResponse<UserAuthResponse>>
    
    @POST("auth/login")
    suspend fun login(@Body request: UserLoginRequest): Response<ApiResponse<UserAuthResponse>>
    
    @GET("auth/me")
    suspend fun getProfile(): Response<ApiResponse<User>>
    
    @POST("auth/logout")
    suspend fun logout(): Response<ApiResponse<Any>>
    
    @POST("auth/change-password")
    suspend fun changePassword(@Body request: ChangePasswordRequest): Response<ApiResponse<Any>>
    
    @DELETE("auth/delete-account")
    suspend fun deleteAccount(): Response<ApiResponse<Any>>
    
    // Shop Discovery
    @GET("shop")
    suspend fun getNearbyShops(
        @Query("lat") latitude: Double,
        @Query("lon") longitude: Double,
        @Query("radius") radius: Double = 5000.0,
        @Query("limit") limit: Int = 20
    ): Response<ApiResponse<List<Shop>>>
    
    // Product Search
    @GET("api/catalog-products/suggest")
    suspend fun searchProducts(
        @Query("keyword") keyword: String,
        @Query("limit") limit: Int = 20
    ): Response<ApiResponse<Map<String, List<CatalogProduct>>>>
}
```

### Network Module Setup

```kotlin
// network/NetworkModule.kt
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import java.util.concurrent.TimeUnit

object NetworkModule {
    
    private const val BASE_URL = "http://10.0.2.2:8080/" // For Android emulator
    // Use your actual server URL for production: "https://your-api.com/"
    
    private val loggingInterceptor = HttpLoggingInterceptor().apply {
        level = HttpLoggingInterceptor.Level.BODY
    }
    
    private val httpClient = OkHttpClient.Builder()
        .addInterceptor(loggingInterceptor)
        .addInterceptor(AuthInterceptor()) // We'll create this
        .connectTimeout(30, TimeUnit.SECONDS)
        .readTimeout(30, TimeUnit.SECONDS)
        .build()
    
    val apiService: ApiService = Retrofit.Builder()
        .baseUrl(BASE_URL)
        .client(httpClient)
        .addConverterFactory(GsonConverterFactory.create())
        .build()
        .create(ApiService::class.java)
}

// network/AuthInterceptor.kt
import okhttp3.Interceptor
import okhttp3.Response

class AuthInterceptor(
    private val tokenProvider: () -> String?
) : Interceptor {
    
    override fun intercept(chain: Interceptor.Chain): Response {
        val originalRequest = chain.request()
        val token = tokenProvider()
        
        return if (token != null) {
            val authenticatedRequest = originalRequest.newBuilder()
                .header("Authorization", "Bearer $token")
                .build()
            chain.proceed(authenticatedRequest)
        } else {
            chain.proceed(originalRequest)
        }
    }
}
```

## Repository Layer

```kotlin
// repository/AuthRepository.kt
class AuthRepository(
    private val apiService: ApiService,
    private val localDataSource: LocalDataSource
) {
    
    suspend fun register(
        name: String,
        email: String,
        password: String
    ): Result<UserAuthResponse> = try {
        val response = apiService.register(UserRegisterRequest(name, email, password))
        
        if (response.isSuccessful && response.body()?.success == true) {
            val authResponse = response.body()!!.data!!
            localDataSource.saveUser(authResponse.user)
            localDataSource.saveToken(authResponse.token)
            Result.success(authResponse)
        } else {
            val errorBody = response.body()
            Result.failure(Exception(errorBody?.message ?: "Registration failed"))
        }
    } catch (e: Exception) {
        Result.failure(e)
    }
    
    suspend fun login(email: String, password: String): Result<UserAuthResponse> = try {
        val response = apiService.login(UserLoginRequest(email, password))
        
        if (response.isSuccessful && response.body()?.success == true) {
            val authResponse = response.body()!!.data!!
            localDataSource.saveUser(authResponse.user)
            localDataSource.saveToken(authResponse.token)
            Result.success(authResponse)
        } else {
            val errorBody = response.body()
            Result.failure(Exception(errorBody?.message ?: "Login failed"))
        }
    } catch (e: Exception) {
        Result.failure(e)
    }
    
    suspend fun logout(): Result<Unit> = try {
        apiService.logout()
        localDataSource.clearUser()
        localDataSource.clearToken()
        Result.success(Unit)
    } catch (e: Exception) {
        // Clear local data even if API call fails
        localDataSource.clearUser()
        localDataSource.clearToken()
        Result.success(Unit)
    }
    
    suspend fun getCurrentUser(): User? = localDataSource.getUser()
    
    fun getToken(): String? = localDataSource.getToken()
    
    fun isLoggedIn(): Boolean = getToken() != null
}

// repository/ShopRepository.kt
class ShopRepository(
    private val apiService: ApiService
) {
    
    suspend fun getNearbyShops(
        latitude: Double,
        longitude: Double,
        radius: Double = 5000.0,
        limit: Int = 20
    ): Result<List<Shop>> = try {
        val response = apiService.getNearbyShops(latitude, longitude, radius, limit)
        
        if (response.isSuccessful && response.body()?.success == true) {
            Result.success(response.body()!!.data ?: emptyList())
        } else {
            Result.failure(Exception(response.body()?.message ?: "Failed to fetch shops"))
        }
    } catch (e: Exception) {
        Result.failure(e)
    }
    
    suspend fun searchProducts(keyword: String, limit: Int = 20): Result<List<CatalogProduct>> = try {
        val response = apiService.searchProducts(keyword, limit)
        
        if (response.isSuccessful && response.body()?.success == true) {
            val products = response.body()!!.data?.get("products") ?: emptyList()
            Result.success(products)
        } else {
            Result.failure(Exception(response.body()?.message ?: "Failed to search products"))
        }
    } catch (e: Exception) {
        Result.failure(e)
    }
}
```

## Local Data Storage

```kotlin
// data/local/LocalDataSource.kt
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map

class LocalDataSource(
    private val dataStore: DataStore<Preferences>
) {
    
    companion object {
        private val USER_TOKEN_KEY = stringPreferencesKey("user_token")
        private val USER_DATA_KEY = stringPreferencesKey("user_data")
    }
    
    suspend fun saveToken(token: String) {
        dataStore.edit { preferences ->
            preferences[USER_TOKEN_KEY] = token
        }
    }
    
    suspend fun saveUser(user: User) {
        dataStore.edit { preferences ->
            preferences[USER_DATA_KEY] = Gson().toJson(user)
        }
    }
    
    fun getToken(): String? {
        return runBlocking {
            dataStore.data.map { preferences ->
                preferences[USER_TOKEN_KEY]
            }.first()
        }
    }
    
    fun getUser(): User? {
        return runBlocking {
            dataStore.data.map { preferences ->
                preferences[USER_DATA_KEY]?.let { userJson ->
                    Gson().fromJson(userJson, User::class.java)
                }
            }.first()
        }
    }
    
    suspend fun clearToken() {
        dataStore.edit { preferences ->
            preferences.remove(USER_TOKEN_KEY)
        }
    }
    
    suspend fun clearUser() {
        dataStore.edit { preferences ->
            preferences.remove(USER_DATA_KEY)
        }
    }
}
```

## ViewModels

```kotlin
// ui/auth/AuthViewModel.kt
class AuthViewModel(
    private val authRepository: AuthRepository
) : ViewModel() {
    
    private val _uiState = MutableStateFlow(AuthUiState())
    val uiState: StateFlow<AuthUiState> = _uiState.asStateFlow()
    
    fun login(email: String, password: String) {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)
            
            authRepository.login(email, password)
                .onSuccess { authResponse ->
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        isLoggedIn = true,
                        user = authResponse.user
                    )
                }
                .onFailure { exception ->
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = exception.message
                    )
                }
        }
    }
    
    fun register(name: String, email: String, password: String) {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)
            
            authRepository.register(name, email, password)
                .onSuccess { authResponse ->
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        isLoggedIn = true,
                        user = authResponse.user
                    )
                }
                .onFailure { exception ->
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = exception.message
                    )
                }
        }
    }
    
    fun logout() {
        viewModelScope.launch {
            authRepository.logout()
            _uiState.value = AuthUiState()
        }
    }
    
    fun checkAuthStatus() {
        val user = authRepository.getCurrentUser()
        val isLoggedIn = authRepository.isLoggedIn()
        _uiState.value = _uiState.value.copy(
            isLoggedIn = isLoggedIn,
            user = user
        )
    }
}

data class AuthUiState(
    val isLoading: Boolean = false,
    val isLoggedIn: Boolean = false,
    val user: User? = null,
    val error: String? = null
)

// ui/shops/ShopsViewModel.kt
class ShopsViewModel(
    private val shopRepository: ShopRepository,
    private val locationProvider: LocationProvider
) : ViewModel() {
    
    private val _uiState = MutableStateFlow(ShopsUiState())
    val uiState: StateFlow<ShopsUiState> = _uiState.asStateFlow()
    
    fun loadNearbyShops() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)
            
            try {
                val location = locationProvider.getCurrentLocation()
                
                shopRepository.getNearbyShops(
                    latitude = location.latitude,
                    longitude = location.longitude
                ).onSuccess { shops ->
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        shops = shops,
                        userLocation = location
                    )
                }.onFailure { exception ->
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = exception.message
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoading = false,
                    error = "Failed to get location: ${e.message}"
                )
            }
        }
    }
    
    fun searchProducts(keyword: String) {
        if (keyword.isBlank()) {
            _uiState.value = _uiState.value.copy(searchResults = emptyList())
            return
        }
        
        viewModelScope.launch {
            shopRepository.searchProducts(keyword)
                .onSuccess { products ->
                    _uiState.value = _uiState.value.copy(searchResults = products)
                }
                .onFailure {
                    _uiState.value = _uiState.value.copy(searchResults = emptyList())
                }
        }
    }
}

data class ShopsUiState(
    val isLoading: Boolean = false,
    val shops: List<Shop> = emptyList(),
    val searchResults: List<CatalogProduct> = emptyList(),
    val userLocation: Location? = null,
    val error: String? = null
)

data class Location(
    val latitude: Double,
    val longitude: Double
)
```

## UI Screens (Jetpack Compose)

```kotlin
// ui/auth/LoginScreen.kt
@Composable
fun LoginScreen(
    viewModel: AuthViewModel,
    onNavigateToRegister: () -> Unit,
    onLoginSuccess: () -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()
    var email by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }
    
    LaunchedEffect(uiState.isLoggedIn) {
        if (uiState.isLoggedIn) {
            onLoginSuccess()
        }
    }
    
    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(16.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Text(
            text = "Shop Near U",
            style = MaterialTheme.typography.headlineLarge,
            modifier = Modifier.padding(bottom = 32.dp)
        )
        
        OutlinedTextField(
            value = email,
            onValueChange = { email = it },
            label = { Text("Email") },
            modifier = Modifier.fillMaxWidth(),
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Email)
        )
        
        Spacer(modifier = Modifier.height(16.dp))
        
        OutlinedTextField(
            value = password,
            onValueChange = { password = it },
            label = { Text("Password") },
            modifier = Modifier.fillMaxWidth(),
            visualTransformation = PasswordVisualTransformation(),
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password)
        )
        
        if (uiState.error != null) {
            Spacer(modifier = Modifier.height(8.dp))
            Text(
                text = uiState.error,
                color = MaterialTheme.colorScheme.error,
                style = MaterialTheme.typography.bodySmall
            )
        }
        
        Spacer(modifier = Modifier.height(24.dp))
        
        Button(
            onClick = { viewModel.login(email, password) },
            modifier = Modifier.fillMaxWidth(),
            enabled = !uiState.isLoading && email.isNotBlank() && password.isNotBlank()
        ) {
            if (uiState.isLoading) {
                CircularProgressIndicator(size = 16.dp)
            } else {
                Text("Login")
            }
        }
        
        Spacer(modifier = Modifier.height(16.dp))
        
        TextButton(onClick = onNavigateToRegister) {
            Text("Don't have an account? Register")
        }
    }
}

// ui/shops/ShopsScreen.kt
@Composable
fun ShopsScreen(
    viewModel: ShopsViewModel,
    onShopClick: (Shop) -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()
    var searchQuery by remember { mutableStateOf("") }
    
    LaunchedEffect(Unit) {
        viewModel.loadNearbyShops()
    }
    
    Column(
        modifier = Modifier.fillMaxSize()
    ) {
        // Search bar
        OutlinedTextField(
            value = searchQuery,
            onValueChange = { query ->
                searchQuery = query
                viewModel.searchProducts(query)
            },
            label = { Text("Search products...") },
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            leadingIcon = {
                Icon(Icons.Default.Search, contentDescription = "Search")
            }
        )
        
        if (uiState.isLoading) {
            Box(
                modifier = Modifier.fillMaxSize(),
                contentAlignment = Alignment.Center
            ) {
                CircularProgressIndicator()
            }
        } else if (uiState.error != null) {
            Box(
                modifier = Modifier.fillMaxSize(),
                contentAlignment = Alignment.Center
            ) {
                Column(horizontalAlignment = Alignment.CenterHorizontally) {
                    Text(
                        text = uiState.error,
                        color = MaterialTheme.colorScheme.error,
                        textAlign = TextAlign.Center
                    )
                    Spacer(modifier = Modifier.height(16.dp))
                    Button(onClick = { viewModel.loadNearbyShops() }) {
                        Text("Retry")
                    }
                }
            }
        } else {
            LazyColumn {
                if (searchQuery.isNotEmpty() && uiState.searchResults.isNotEmpty()) {
                    item {
                        Text(
                            text = "Search Results",
                            style = MaterialTheme.typography.headlineSmall,
                            modifier = Modifier.padding(16.dp)
                        )
                    }
                    items(uiState.searchResults) { product ->
                        ProductItem(product = product)
                    }
                    item {
                        Divider(modifier = Modifier.padding(vertical = 16.dp))
                        Text(
                            text = "Nearby Shops",
                            style = MaterialTheme.typography.headlineSmall,
                            modifier = Modifier.padding(horizontal = 16.dp)
                        )
                    }
                }
                
                items(uiState.shops) { shop ->
                    ShopItem(
                        shop = shop,
                        onClick = { onShopClick(shop) }
                    )
                }
            }
        }
    }
}

@Composable
fun ShopItem(
    shop: Shop,
    onClick: () -> Unit
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp, vertical = 4.dp)
            .clickable { onClick() },
        elevation = CardDefaults.cardElevation(defaultElevation = 4.dp)
    ) {
        Column(
            modifier = Modifier.padding(16.dp)
        ) {
            Text(
                text = shop.name,
                style = MaterialTheme.typography.headlineSmall
            )
            Text(
                text = "Owner: ${shop.ownerName}",
                style = MaterialTheme.typography.bodyMedium,
                color = MaterialTheme.colorScheme.onSurfaceVariant
            )
            Text(
                text = shop.address,
                style = MaterialTheme.typography.bodySmall,
                color = MaterialTheme.colorScheme.onSurfaceVariant
            )
            
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(top = 8.dp),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text(
                    text = shop.type,
                    style = MaterialTheme.typography.labelMedium,
                    color = MaterialTheme.colorScheme.primary
                )
                
                shop.distance?.let { distance ->
                    Text(
                        text = "${String.format("%.1f", distance)}m away",
                        style = MaterialTheme.typography.labelSmall,
                        color = MaterialTheme.colorScheme.outline
                    )
                }
            }
            
            if (shop.supportsDelivery) {
                Row(
                    modifier = Modifier.padding(top = 4.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(
                        Icons.Default.LocalShipping,
                        contentDescription = "Delivery",
                        tint = MaterialTheme.colorScheme.primary,
                        modifier = Modifier.size(16.dp)
                    )
                    Spacer(modifier = Modifier.width(4.dp))
                    Text(
                        text = "Delivery Available",
                        style = MaterialTheme.typography.labelSmall,
                        color = MaterialTheme.colorScheme.primary
                    )
                }
            }
        }
    }
}

@Composable
fun ProductItem(product: CatalogProduct) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp, vertical = 4.dp),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Row(
            modifier = Modifier.padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            AsyncImage(
                model = product.imageUrl,
                contentDescription = product.name,
                modifier = Modifier
                    .size(60.dp)
                    .clip(RoundedCornerShape(8.dp)),
                placeholder = painterResource(R.drawable.placeholder_product),
                error = painterResource(R.drawable.placeholder_product)
            )
            
            Spacer(modifier = Modifier.width(16.dp))
            
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = product.name,
                    style = MaterialTheme.typography.titleMedium
                )
                product.brand?.let { brand ->
                    Text(
                        text = brand,
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                }
                Text(
                    text = product.category,
                    style = MaterialTheme.typography.labelSmall,
                    color = MaterialTheme.colorScheme.primary
                )
            }
        }
    }
}
```

## Location Services

```kotlin
// location/LocationProvider.kt
class LocationProvider(
    private val context: Context,
    private val fusedLocationClient: FusedLocationProviderClient
) {
    
    suspend fun getCurrentLocation(): Location = suspendCancellableCoroutine { continuation ->
        if (ActivityCompat.checkSelfPermission(
                context,
                Manifest.permission.ACCESS_FINE_LOCATION
            ) != PackageManager.PERMISSION_GRANTED
        ) {
            continuation.resumeWithException(SecurityException("Location permission not granted"))
            return@suspendCancellableCoroutine
        }
        
        fusedLocationClient.lastLocation
            .addOnSuccessListener { androidLocation ->
                if (androidLocation != null) {
                    continuation.resume(
                        Location(
                            latitude = androidLocation.latitude,
                            longitude = androidLocation.longitude
                        )
                    )
                } else {
                    continuation.resumeWithException(Exception("Unable to get location"))
                }
            }
            .addOnFailureListener { exception ->
                continuation.resumeWithException(exception)
            }
    }
}
```

## Navigation

```kotlin
// navigation/NavGraph.kt
@Composable
fun NavGraph(
    navController: NavHostController,
    startDestination: String
) {
    NavHost(
        navController = navController,
        startDestination = startDestination
    ) {
        composable("login") {
            val authViewModel: AuthViewModel = koinViewModel()
            LoginScreen(
                viewModel = authViewModel,
                onNavigateToRegister = {
                    navController.navigate("register")
                },
                onLoginSuccess = {
                    navController.navigate("shops") {
                        popUpTo("login") { inclusive = true }
                    }
                }
            )
        }
        
        composable("register") {
            val authViewModel: AuthViewModel = koinViewModel()
            RegisterScreen(
                viewModel = authViewModel,
                onNavigateToLogin = {
                    navController.popBackStack()
                },
                onRegisterSuccess = {
                    navController.navigate("shops") {
                        popUpTo("register") { inclusive = true }
                    }
                }
            )
        }
        
        composable("shops") {
            val shopsViewModel: ShopsViewModel = koinViewModel()
            ShopsScreen(
                viewModel = shopsViewModel,
                onShopClick = { shop ->
                    navController.navigate("shop_detail/${shop.id}")
                }
            )
        }
        
        composable(
            "shop_detail/{shopId}",
            arguments = listOf(navArgument("shopId") { type = NavType.IntType })
        ) { backStackEntry ->
            val shopId = backStackEntry.arguments?.getInt("shopId") ?: return@composable
            // ShopDetailScreen(shopId = shopId)
        }
    }
}
```

## MainActivity

```kotlin
// MainActivity.kt
class MainActivity : ComponentActivity() {
    
    private val authViewModel: AuthViewModel by viewModel()
    
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        
        // Check location permissions
        requestLocationPermissions()
        
        setContent {
            ShopNearUTheme {
                val navController = rememberNavController()
                val uiState by authViewModel.uiState.collectAsState()
                
                LaunchedEffect(Unit) {
                    authViewModel.checkAuthStatus()
                }
                
                NavGraph(
                    navController = navController,
                    startDestination = if (uiState.isLoggedIn) "shops" else "login"
                )
            }
        }
    }
    
    private fun requestLocationPermissions() {
        if (ActivityCompat.checkSelfPermission(
                this,
                Manifest.permission.ACCESS_FINE_LOCATION
            ) != PackageManager.PERMISSION_GRANTED
        ) {
            ActivityCompat.requestPermissions(
                this,
                arrayOf(Manifest.permission.ACCESS_FINE_LOCATION),
                LOCATION_PERMISSION_REQUEST_CODE
            )
        }
    }
    
    companion object {
        private const val LOCATION_PERMISSION_REQUEST_CODE = 1
    }
}
```

## Dependency Injection Setup

```kotlin
// di/AppModule.kt
val appModule = module {
    
    // DataStore
    single<DataStore<Preferences>> {
        get<Context>().dataStore
    }
    
    // Local Data Source
    single { LocalDataSource(get()) }
    
    // Location
    single { LocationServices.getFusedLocationProviderClient(get<Context>()) }
    single { LocationProvider(get(), get()) }
    
    // Network
    single { NetworkModule.apiService }
    
    // Repositories
    single { AuthRepository(get(), get()) }
    single { ShopRepository(get()) }
    
    // ViewModels
    viewModel { AuthViewModel(get()) }
    viewModel { ShopsViewModel(get(), get()) }
}

// Application.kt
class ShopNearUApplication : Application() {
    
    override fun onCreate() {
        super.onCreate()
        
        startKoin {
            androidContext(this@ShopNearUApplication)
            modules(appModule)
        }
    }
}

val Context.dataStore: DataStore<Preferences> by preferencesDataStore(name = "settings")
```

## Permissions (AndroidManifest.xml)

```xml
<uses-permission android:name="android.permission.INTERNET" />
<uses-permission android:name="android.permission.ACCESS_FINE_LOCATION" />
<uses-permission android:name="android.permission.ACCESS_COARSE_LOCATION" />
```

## Key Features Implementation Notes

1. **Authentication**: Token-based auth with local storage
2. **Location Services**: GPS-based shop discovery
3. **Offline Support**: Cache user data and last known shops
4. **Search**: Real-time product search with debouncing
5. **Error Handling**: Comprehensive error states and retry mechanisms
6. **UI/UX**: Material Design 3 with Compose
7. **Architecture**: MVVM + Repository pattern with clean separation

## Next Steps

1. Implement the remaining screens (RegisterScreen, ShopDetailScreen)
2. Add proper error handling and loading states
3. Implement offline caching with Room database
4. Add maps integration for visual shop locations
5. Implement push notifications for shop updates
6. Add product ordering/cart functionality (requires backend expansion)
7. Implement proper CI/CD with automated testing

This specification provides a solid foundation for building a robust Android app that integrates with your shop-near-u backend API.